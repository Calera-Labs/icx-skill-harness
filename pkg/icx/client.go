package icx

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	statusLocalFallback = "LOCAL_FALLBACK"
	localStoreLocation  = "process-local store"
)

// Client represents the ICX REST Gateway Client
type Client struct {
	BaseURL       string
	APIKey        string
	SpaceID       string
	BYOKProvider  string
	BYOKKey       string
	BYOKModel     string
	HTTPClient    *http.Client
	mu            sync.RWMutex
	localFallback bool
	localStore    map[string][]MemoryRegisterMatch
}

// Config holds configuration parameters for the ICX Client
type Config struct {
	BaseURL       string
	APIKey        string
	SpaceID       string
	BYOKProvider  string
	BYOKKey       string
	BYOKModel     string
	Timeout       time.Duration
	LocalFallback bool // process-local store only; skip hosted ICX even if APIKey is set
}

// NewClient creates an ICX client with configuration
func NewClient(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://icx.api.caleralabs.com/v1"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.BYOKModel == "" {
		cfg.BYOKModel = "gemini-3.5-flash-lite"
	}
	if cfg.BYOKProvider == "" {
		cfg.BYOKProvider = "gemini"
	}

	return &Client{
		BaseURL:       strings.TrimRight(cfg.BaseURL, "/"),
		APIKey:        cfg.APIKey,
		SpaceID:       cfg.SpaceID,
		BYOKProvider:  cfg.BYOKProvider,
		BYOKKey:       cfg.BYOKKey,
		BYOKModel:     cfg.BYOKModel,
		localFallback: cfg.LocalFallback || strings.TrimSpace(cfg.APIKey) == "",
		HTTPClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		localStore: make(map[string][]MemoryRegisterMatch),
	}
}

func (c *Client) buildHeaders(spaceID string) http.Header {
	targetSpace := spaceID
	if targetSpace == "" {
		targetSpace = c.SpaceID
	}

	h := make(http.Header)
	h.Set("Authorization", "Bearer "+c.APIKey)
	h.Set("Content-Type", "application/json")
	if targetSpace != "" {
		h.Set("X-Space-ID", targetSpace)
	}

	// Never forward the operator's model API key to Calera. Provider/model
	// names are non-secret hints only.
	if c.BYOKProvider != "" {
		h.Set("X-LLM-Provider", c.BYOKProvider)
	}
	if c.BYOKModel != "" {
		h.Set("X-LLM-Model", c.BYOKModel)
	}
	return h
}

func contentHash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

func (c *Client) resolveSpace(spaceID string) string {
	if strings.TrimSpace(spaceID) != "" {
		return spaceID
	}
	return c.SpaceID
}

func truncateBody(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 240 {
		return s[:240] + "…"
	}
	return s
}

func (c *Client) storeLocal(req IngestTextRequest, elapsedMs float64) *IngestResponse {
	space := c.resolveSpace(req.SpaceID)
	hash := contentHash(req.Text)

	c.mu.Lock()
	c.localStore[space] = append(c.localStore[space], MemoryRegisterMatch{
		Content:       req.Text,
		DocumentTitle: req.Filename,
		Location:      localStoreLocation,
		Confidence:    0,
		Subject:       req.Family,
	})
	c.mu.Unlock()

	return &IngestResponse{
		Success:       true,
		Status:        statusLocalFallback,
		SpaceID:       space,
		Filename:      req.Filename,
		ContentHash:   hash,
		MerkleHash:    hash, // alias: SHA-256 of the payload, not a Merkle tree
		LatencyMs:     elapsedMs,
		LocalFallback: true,
	}
}

func (c *Client) recallLocal(req RecallRequest, elapsedMs float64) *RecallResponse {
	space := c.resolveSpace(req.SpaceID)

	c.mu.RLock()
	registers := c.localStore[space]
	c.mu.RUnlock()

	matches := make([]MemoryRegisterMatch, 0)
	qLower := strings.ToLower(req.Query)
	qWords := strings.Fields(qLower)

	for _, reg := range registers {
		cLower := strings.ToLower(reg.Content)
		score := 0
		for _, w := range qWords {
			if len(w) > 2 && strings.Contains(cLower, w) {
				score++
			}
		}
		if score > 0 {
			regMatch := reg
			regMatch.Location = localStoreLocation
			regMatch.Confidence = float64(score) / float64(len(qWords)+1)
			matches = append(matches, regMatch)
		}
	}

	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}
	if len(matches) > topK {
		matches = matches[:topK]
	}

	return &RecallResponse{
		Success:       true,
		SpaceID:       space,
		Query:         req.Query,
		Matches:       matches,
		LatencyMs:     elapsedMs,
		LocalFallback: true,
	}
}

// IngestText ingests raw text or skill definitions into ICX, or the process-local store.
func (c *Client) IngestText(req IngestTextRequest) (*IngestResponse, error) {
	return c.IngestTextWithContext(context.Background(), req)
}

// IngestTextWithContext ingests raw text with explicit context cancellation.
// With no API key (or LocalFallback), this writes to a process-local map and
// labels the result LOCAL_FALLBACK. With an API key, HTTP errors are returned
// to the caller — they are not reported as a successful ICX ingest.
func (c *Client) IngestTextWithContext(ctx context.Context, req IngestTextRequest) (*IngestResponse, error) {
	t0 := time.Now()
	if c.localFallback {
		return c.storeLocal(req, float64(time.Since(t0).Microseconds())/1000.0), nil
	}

	url := fmt.Sprintf("%s/ingest/text", c.BaseURL)
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ingest request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header = c.buildHeaders(req.SpaceID)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("icx ingest request failed: %w", err)
	}
	defer resp.Body.Close()
	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("icx ingest HTTP %d: %s", resp.StatusCode, truncateBody(respBytes))
	}

	var res IngestResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to decode ingest response: %w", err)
	}
	if res.ContentHash == "" {
		res.ContentHash = res.MerkleHash
	}
	res.LatencyMs = float64(time.Since(t0).Microseconds()) / 1000.0
	return &res, nil
}

// IngestBatch ingests multiple text items into the ICX Volumetric Lattice
func (c *Client) IngestBatch(req IngestBatchRequest) (*IngestBatchResponse, error) {
	return c.IngestBatchWithContext(context.Background(), req)
}

// IngestBatchWithContext ingests multiple text items with context
func (c *Client) IngestBatchWithContext(ctx context.Context, req IngestBatchRequest) (*IngestBatchResponse, error) {
	t0 := time.Now()
	totalFacts := 0
	h := sha256.New()
	var errors []string
	anyLocal := false

	for _, item := range req.Items {
		if item.SpaceID == "" {
			item.SpaceID = req.SpaceID
		}
		res, err := c.IngestTextWithContext(ctx, item)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", item.Filename, err))
			continue
		}
		if res.LocalFallback {
			anyLocal = true
		}
		totalFacts += res.FactsCrystallized
		hash := res.ContentHash
		if hash == "" {
			hash = res.MerkleHash
		}
		h.Write([]byte(hash))
	}

	elapsed := float64(time.Since(t0).Microseconds()) / 1000.0
	combined := hex.EncodeToString(h.Sum(nil))
	status := ""
	if anyLocal && len(errors) == 0 {
		status = statusLocalFallback
	}
	return &IngestBatchResponse{
		Success:           len(errors) == 0,
		Status:            status,
		TotalItems:        len(req.Items),
		FactsCrystallized: totalFacts,
		ContentHash:       combined,
		MerkleHash:        combined,
		LatencyMs:         elapsed,
		LocalFallback:     anyLocal,
		Errors:            errors,
	}, nil
}

// Recall queries the ICX memory space for matching registers
func (c *Client) Recall(req RecallRequest) (*RecallResponse, error) {
	return c.RecallWithContext(context.Background(), req)
}

// RecallWithContext queries ICX (or the process-local store) with context.
// Hosted failures are returned as errors when an API key is configured.
func (c *Client) RecallWithContext(ctx context.Context, req RecallRequest) (*RecallResponse, error) {
	t0 := time.Now()
	if c.localFallback {
		return c.recallLocal(req, float64(time.Since(t0).Microseconds())/1000.0), nil
	}

	url := fmt.Sprintf("%s/recall", c.BaseURL)
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recall request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header = c.buildHeaders(req.SpaceID)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("icx recall request failed: %w", err)
	}
	defer resp.Body.Close()
	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("icx recall HTTP %d: %s", resp.StatusCode, truncateBody(respBytes))
	}

	var res RecallResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to decode recall response: %w", err)
	}
	res.LatencyMs = float64(time.Since(t0).Microseconds()) / 1000.0
	return &res, nil
}

// ChatCompletion executes a chat completion via the ICX Gateway with BYOK reasoning
func (c *Client) ChatCompletion(req ChatCompletionRequest, spaceID string) (*ChatCompletionResponse, error) {
	return c.ChatCompletionWithContext(context.Background(), req, spaceID)
}

// ChatCompletionWithContext executes chat completion with context
func (c *Client) ChatCompletionWithContext(ctx context.Context, req ChatCompletionRequest, spaceID string) (*ChatCompletionResponse, error) {
	if c.localFallback {
		return nil, fmt.Errorf("icx chat is not available in local-fallback mode")
	}

	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header = c.buildHeaders(spaceID)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request to ICX chat completion failed: %w", err)
	}
	defer resp.Body.Close()
	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("icx chat HTTP %d: %s", resp.StatusCode, truncateBody(respBytes))
	}

	var res ChatCompletionResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to parse chat response: %w", err)
	}
	return &res, nil
}
