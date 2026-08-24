package gateway

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/caleralabs/icx-skill-harness/pkg/byok"
	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/router"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

// Server implements the OpenAI-compatible HTTP Gateway Daemon with JIT Lattice Routing
type Server struct {
	config       GatewayConfig
	baseRegistry *skills.SkillRegistry
	icxClient    *icx.Client
	llmClient    byok.LLMClient
	router       *router.LatticeSkillRouter
	stats        GatewayStats
	httpServer   *http.Server
	listener     net.Listener
}

// NewServer initializes a new gateway daemon instance
func NewServer(cfg GatewayConfig, baseReg *skills.SkillRegistry, icxClient *icx.Client, llm byok.LLMClient) *Server {
	if baseReg == nil {
		baseReg = skills.NewSkillRegistry()
	}

	routerCfg := router.DefaultRouterConfig()
	if cfg.MaxToolsPerViewport > 0 {
		routerCfg.MaxToolsPerViewport = cfg.MaxToolsPerViewport
	}
	jitRouter := router.NewLatticeSkillRouter(baseReg, icxClient, routerCfg)

	s := &Server{
		config:       cfg,
		baseRegistry: baseReg,
		icxClient:    icxClient,
		llmClient:    llm,
		router:       jitRouter,
	}
	s.stats.StartTime = time.Now()
	return s
}

// Start launches the HTTP proxy server
func (s *Server) Start(ctx context.Context) error {
	if strings.TrimSpace(s.config.GatewayToken) == "" {
		return fmt.Errorf("gateway token is required; set ICX_GATEWAY_TOKEN or pass --gateway-token")
	}
	host := s.config.Host
	if host == "" {
		host = "127.0.0.1"
		s.config.Host = host
	}
	if !IsLoopbackHost(host) && !s.config.AllowPublicBind {
		return fmt.Errorf("refusing to bind non-loopback address %q (pass --allow-public-bind)", host)
	}
	if s.config.MaxBodyBytes <= 0 {
		s.config.MaxBodyBytes = defaultMaxBodyBytes
	}
	if s.config.MaxIncomingTools <= 0 {
		s.config.MaxIncomingTools = defaultMaxTools
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/v1/chat/completions", s.requireAuth(s.handleChatCompletions))
	mux.HandleFunc("/chat/completions", s.requireAuth(s.handleChatCompletions))
	mux.HandleFunc("/v1/models", s.requireAuth(s.handleModels))
	mux.HandleFunc("/models", s.requireAuth(s.handleModels))
	mux.HandleFunc("/v1/skills/register", s.requireAuth(s.handleSkillRegister))
	mux.HandleFunc("/v1/stats", s.requireAuth(s.handleStats))
	mux.HandleFunc("/v1/health", s.requireAuth(s.handleHealth))

	addr := fmt.Sprintf("%s:%d", host, s.config.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.corsMiddleware(mux),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind address %s: %w", addr, err)
	}
	s.listener = ln

	fmt.Printf("\n========================================================================================\n")
	fmt.Printf("          🚀 CALERA ICX UNIVERSAL SKILL GATEWAY DAEMON (OPENAI-COMPATIBLE)             \n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" Listening Endpoint             : http://%s\n", ln.Addr().String())
	fmt.Printf(" OpenAI Drop-In URL             : http://%s/v1/chat/completions\n", ln.Addr().String())
	fmt.Printf(" Active ICX Memory Space        : %s\n", s.config.SpaceID)
	fmt.Printf(" BYOK Upstream Model            : %s (%s)\n", s.config.DefaultModel, s.config.DefaultProvider)
	fmt.Printf(" Pre-Loaded Base Skills         : %d Skills (%d Tools)\n", s.baseRegistry.Count(), len(s.baseRegistry.GetAllTools()))
	fmt.Printf(" JIT Micro-Viewport Capacity    : Max %d tools / request\n", s.config.MaxToolsPerViewport)
	fmt.Printf(" State Crystallization Engine   : %v\n", s.config.EnableCrystallization)
	fmt.Printf(" Auth                           : bearer gateway token required (except /healthz)\n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" This daemon holds your model key. Keep it on localhost. Pass Authorization: Bearer <gateway-token>.\n\n")

	go func() {
		if err := s.httpServer.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("Gateway server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the gateway server
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// Addr returns the bound address of the server
func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
}

// handleChatCompletions is the drop-in OpenAI-compatible chat completions handler
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !tokenEqual(extractBearer(r), s.config.GatewayToken) {
			writeJSONError(w, http.StatusUnauthorized, "missing or invalid gateway token")
			return
		}
		next(w, r)
	}
}

func (s *Server) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	s.stats.TotalRequests.Add(1)

	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Only POST method is supported")
		return
	}

	limit := s.config.MaxBodyBytes
	if limit <= 0 {
		limit = defaultMaxBodyBytes
	}
	r.Body = http.MaxBytesReader(w, r.Body, limit)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": {"message": "Failed to read request body", "type": "invalid_request_error"}}`, http.StatusBadRequest)
		return
	}

	var req OpenAIChatCompletionRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": {"message": "Invalid JSON: %s", "type": "invalid_request_error"}}`, err.Error()), http.StatusBadRequest)
		return
	}

	// 1. Extract full conversational context & latest user prompt
	var userPrompt string
	var systemInstruction string
	var contents []byok.GeminiContent

	for _, msg := range req.Messages {
		role := strings.ToLower(msg.Role)
		if role == "system" {
			systemInstruction = msg.Content
		}
		if role == "user" {
			userPrompt = msg.Content
		}

		byokRole := "user"
		if role == "assistant" {
			byokRole = "model"
		}

		var parts []byok.GeminiPart
		if msg.Content != "" {
			parts = append(parts, byok.GeminiPart{Text: msg.Content})
		}
		for _, tc := range msg.ToolCalls {
			var argsMap map[string]any
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &argsMap)
			parts = append(parts, byok.GeminiPart{
				FunctionCall: &byok.GeminiFunctionCall{
					Name: tc.Function.Name,
					Args: argsMap,
				},
			})
		}

		if len(parts) > 0 {
			contents = append(contents, byok.GeminiContent{
				Role:  byokRole,
				Parts: parts,
			})
		}
	}

	// 2. Bring Your Own Skills (BYOS): Build dynamic request registry
	reqRegistry := s.baseRegistry
	totalIncomingTools := 0
	monolithicTokens := 0

	if len(req.Tools) > 0 {
		maxTools := s.config.MaxIncomingTools
		if maxTools <= 0 {
			maxTools = defaultMaxTools
		}
		if len(req.Tools) > maxTools {
			writeJSONError(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("too many tools (max %d)", maxTools))
			return
		}
		reqRegistry = skills.NewSkillRegistry()
		totalIncomingTools = len(req.Tools)

		for _, openAITool := range req.Tools {
			toolDef := skills.ConvertOpenAIToolToDefinition(
				openAITool.Function.Name,
				openAITool.Function.Description,
				openAITool.Function.Parameters,
			)
			skill := skills.CreateSkillFromTool(toolDef)
			reqRegistry.Register(skill)
		}
		monolithicTokens = reqRegistry.TotalMonolithicTokens()
	} else {
		totalIncomingTools = len(s.baseRegistry.GetAllTools())
		monolithicTokens = s.baseRegistry.TotalMonolithicTokens()
	}

	// 3. JIT Micro-Viewport Routing
	var activeTools []skills.ToolDefinition
	var viewport *skills.SkillViewport
	var routingDuration time.Duration

	if reqRegistry.Count() > 0 && !req.ICXDisablePrune {
		rCfg := router.DefaultRouterConfig()
		rCfg.MaxToolsPerViewport = s.config.MaxToolsPerViewport
		reqRouter := router.NewLatticeSkillRouter(reqRegistry, s.icxClient, rCfg)

		tRoute := time.Now()
		vp, err := reqRouter.RouteQueryWithContext(r.Context(), userPrompt)
		routingDuration = time.Since(tRoute)

		if err == nil && vp != nil {
			viewport = vp
			if !vp.IsRefusal {
				activeTools = vp.ActiveTools
			}
		}
	}

	prunedCount := totalIncomingTools - len(activeTools)
	if prunedCount < 0 {
		prunedCount = 0
	}
	s.stats.TotalToolsPruned.Add(int64(prunedCount))

	// 4. Upstream LLM Execution
	effectiveModel := req.Model
	if effectiveModel == "" {
		effectiveModel = s.config.DefaultModel
	}

	modelResp, err := s.llmClient.GenerateContentWithContext(r.Context(), contents, activeTools, systemInstruction)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "Upstream LLM error")
		return
	}

	// 5. State Crystallization in Calera ICX Space
	merkleSeal := ""
	effectiveSpace := req.ICXSpaceID
	if effectiveSpace == "" {
		effectiveSpace = s.config.SpaceID
	}

	if s.config.EnableCrystallization && s.icxClient != nil && (modelResp.ToolCall != nil || len(modelResp.TextResponse) > 250) {
		ingestContent := modelResp.TextResponse
		if modelResp.ToolCall != nil {
			callJSON, _ := json.Marshal(modelResp.ToolCall)
			ingestContent = string(callJSON)
		}

		ingestResp, err := s.icxClient.IngestTextWithContext(r.Context(), icx.IngestTextRequest{
			SpaceID:  effectiveSpace,
			Text:     ingestContent,
			Filename: "gateway_turn.json",
			Family:   "state.register",
		})
		if err == nil && ingestResp != nil && !ingestResp.LocalFallback {
			merkleSeal = ingestResp.ContentHash
			if merkleSeal == "" {
				merkleSeal = ingestResp.MerkleHash
			}
			s.stats.TotalStatesStored.Add(1)
		}
	}

	// 6. Token & Dollar Economics Calculation
	actualPromptTokens := modelResp.PromptTokens
	if actualPromptTokens == 0 {
		// Fallback approximation
		actualPromptTokens = len(userPrompt)/4 + (len(activeTools) * 110) + 40
	}
	completionTokens := modelResp.ResponseTokens
	if completionTokens == 0 {
		completionTokens = len(modelResp.TextResponse)/4 + 20
	}

	monoEstimatedTokens := monolithicTokens + len(userPrompt)/4 + 80
	if monoEstimatedTokens < actualPromptTokens {
		monoEstimatedTokens = actualPromptTokens
	}

	tokensSaved := monoEstimatedTokens - actualPromptTokens
	if tokensSaved < 0 {
		tokensSaved = 0
	}

	savingsPct := 0.0
	if monoEstimatedTokens > 0 {
		savingsPct = (float64(tokensSaved) / float64(monoEstimatedTokens)) * 100.0
	}

	s.stats.TotalTokensSaved.Add(int64(tokensSaved))
	s.stats.TotalTokensServed.Add(int64(actualPromptTokens + completionTokens))

	// 7. Format Standard OpenAI Response
	completionID := "chatcmpl-" + generateRandomHex(12)
	createdTimestamp := time.Now().Unix()

	choice := OpenAIChoice{
		Index:        0,
		FinishReason: "stop",
		Message: OpenAIMessage{
			Role: "assistant",
		},
	}

	if modelResp.ToolCall != nil {
		choice.FinishReason = "tool_calls"
		argsJSON, _ := json.Marshal(modelResp.ToolCall.Args)
		choice.Message.ToolCalls = []OpenAIToolCall{
			{
				ID:   "call_" + generateRandomHex(8),
				Type: "function",
				Function: OpenAIFunctionCallData{
					Name:      modelResp.ToolCall.Name,
					Arguments: string(argsJSON),
				},
			},
		}
	} else {
		choice.Message.Content = modelResp.TextResponse
		if viewport != nil && viewport.IsRefusal {
			choice.Message.Content = fmt.Sprintf("[CALERA_ICX_SAFE_REFUSAL]: %s", viewport.RefusalReason)
		}
	}

	openAIResp := OpenAIChatCompletionResponse{
		ID:                completionID,
		Object:            "chat.completion",
		Created:           createdTimestamp,
		Model:             effectiveModel,
		SystemFingerprint: "fp_calera_icx_" + generateRandomHex(4),
		Choices:           []OpenAIChoice{choice},
		Usage: OpenAIUsage{
			PromptTokens:     actualPromptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      actualPromptTokens + completionTokens,
			MonolithicTokens: monoEstimatedTokens,
			ICXTokensSaved:   tokensSaved,
			ICXSavingsPct:    savingsPct,
			ICXMerkleSeal:    merkleSeal,
			ICXSpaceID:       effectiveSpace,
		},
	}

	// Response Headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-ICX-Tokens-Saved", fmt.Sprintf("%d", tokensSaved))
	w.Header().Set("X-ICX-Savings-Pct", fmt.Sprintf("%.2f%%", savingsPct))
	w.Header().Set("X-ICX-Space-ID", effectiveSpace)
	if merkleSeal != "" {
		w.Header().Set("X-ICX-Merkle-Seal", merkleSeal)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(openAIResp)

	totalLatency := time.Since(t0)
	dollarSaved := (float64(tokensSaved) / 1000000.0) * 5.0 // Standard $5/1M token benchmark

	// Log live telemetry to console
	fmt.Printf("⚡ [ICX PROXY] %s | %s | Saved %d tokens (%.1f%%) | Latency: %.1fms (Router: %.1fms) | Space: %s\n",
		time.Now().Format("15:04:05"), effectiveModel, tokensSaved, savingsPct,
		float64(totalLatency.Microseconds())/1000.0, float64(routingDuration.Microseconds())/1000.0, effectiveSpace)
	if len(activeTools) > 0 {
		var toolNames []string
		for _, at := range activeTools {
			toolNames = append(toolNames, at.Name)
		}
		fmt.Printf("   ↳ Injected Viewport: [%s] (Pruned %d tools | Saved ~$%.4f)\n",
			strings.Join(toolNames, ", "), prunedCount, dollarSaved)
	}
	if merkleSeal != "" {
		fmt.Printf("   ↳ Crystallized State Seal: %s\n", merkleSeal[:min(16, len(merkleSeal))])
	}
}

// handleModels returns available models
func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	resp := OpenAIModelListResponse{
		Object: "list",
		Data: []OpenAIModelItem{
			{ID: "gemini-3.5-flash-lite", Object: "model", Created: 1700000000, OwnedBy: "google"},
			{ID: "gemini-3.6-flash", Object: "model", Created: 1700000000, OwnedBy: "google"},
			{ID: "gpt-4o", Object: "model", Created: 1700000000, OwnedBy: "openai"},
			{ID: "gpt-4o-mini", Object: "model", Created: 1700000000, OwnedBy: "openai"},
			{ID: "deepseek-chat", Object: "model", Created: 1700000000, OwnedBy: "deepseek"},
			{ID: "claude-3-5-sonnet", Object: "model", Created: 1700000000, OwnedBy: "anthropic"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// handleStats returns live aggregate metrics
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	reqs := s.stats.TotalRequests.Load()
	saved := s.stats.TotalTokensSaved.Load()
	served := s.stats.TotalTokensServed.Load()
	pruned := s.stats.TotalToolsPruned.Load()
	states := s.stats.TotalStatesStored.Load()

	totalPotential := saved + served
	savingsPct := 0.0
	if totalPotential > 0 {
		savingsPct = (float64(saved) / float64(totalPotential)) * 100.0
	}

	snap := StatsSnapshot{
		TotalRequests:        reqs,
		TotalToolsPruned:     pruned,
		TotalTokensSaved:     saved,
		TotalTokensServed:    served,
		EstimatedDollarSaved: (float64(saved) / 1000000.0) * 5.0,
		TotalStatesStored:    states,
		OverallSavingsPct:    savingsPct,
		UptimeSeconds:        int64(time.Since(s.stats.StartTime).Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snap)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":      "HEALTHY",
		"version":     "1.0.0",
		"service":     "calera-icx-gateway",
		"base_skills": s.baseRegistry.Count(),
	})
}

// handleSkillRegister dynamically registers a skill via HTTP POST
func (s *Server) handleSkillRegister(w http.ResponseWriter, r *http.Request) {
	if !s.config.EnableHotRegister {
		writeJSONError(w, http.StatusNotFound, "hot skill registration is disabled")
		return
	}
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	var toolDef skills.ToolDefinition
	if err := json.NewDecoder(r.Body).Decode(&toolDef); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	skill := skills.CreateSkillFromTool(toolDef)
	s.baseRegistry.Register(skill)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":       "REGISTERED",
		"skill_id":     skill.ID,
		"tool_name":    toolDef.Name,
		"merkle_seal":  skill.MerkleSeal,
		"total_skills": s.baseRegistry.Count(),
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isLocalhostOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-ICX-Gateway-Token, X-ICX-Space-ID, X-ICX-Session-ID")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func generateRandomHex(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
