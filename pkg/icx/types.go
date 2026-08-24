package icx

import "time"

// IngestTextRequest payload for POST /v1/ingest/text
type IngestTextRequest struct {
	Text     string `json:"text"`
	Filename string `json:"filename,omitempty"`
	SpaceID  string `json:"space_id,omitempty"`
	Family   string `json:"family,omitempty"`
}

// IngestResponse represents the response from ICX ingestion.
// MerkleHash is a SHA-256 of the payload (kept for compatibility); prefer ContentHash.
// LocalFallback is true when nothing was written to hosted ICX.
type IngestResponse struct {
	Success           bool    `json:"success"`
	Status            string  `json:"status,omitempty"`
	SpaceID           string  `json:"space_id,omitempty"`
	Filename          string  `json:"filename,omitempty"`
	FactsCrystallized int     `json:"facts_crystallized,omitempty"`
	ContentHash       string  `json:"content_hash,omitempty"`
	MerkleHash        string  `json:"merkle_hash,omitempty"`
	LatencyMs         float64 `json:"latency_ms,omitempty"`
	LocalFallback     bool    `json:"local_fallback,omitempty"`
	Error             string  `json:"error,omitempty"`
}

// IngestBatchRequest payload for batch ingestion
type IngestBatchRequest struct {
	Items   []IngestTextRequest `json:"items"`
	SpaceID string              `json:"space_id,omitempty"`
}

// IngestBatchResponse represents batch ingestion results
type IngestBatchResponse struct {
	Success           bool     `json:"success"`
	Status            string   `json:"status,omitempty"`
	TotalItems        int      `json:"total_items"`
	FactsCrystallized int      `json:"facts_crystallized"`
	ContentHash       string   `json:"content_hash,omitempty"`
	MerkleHash        string   `json:"merkle_hash"`
	LatencyMs         float64  `json:"latency_ms"`
	LocalFallback     bool     `json:"local_fallback,omitempty"`
	Errors            []string `json:"errors,omitempty"`
}

// RecallRequest payload for POST /v1/recall or /v1/symbolic/search
type RecallRequest struct {
	Query   string `json:"query"`
	SpaceID string `json:"space_id,omitempty"`
	TopK    int    `json:"top_k,omitempty"`
	Family  string `json:"family,omitempty"`
}

// MemoryRegisterMatch represents a single crystallized memory register from ICX
type MemoryRegisterMatch struct {
	Content       string  `json:"content"`
	DocumentTitle string  `json:"document_title"`
	Location      string  `json:"location"`
	Confidence    float64 `json:"confidence"`
	Subject       string  `json:"subject,omitempty"`
	Relation      string  `json:"relation,omitempty"`
}

// RecallResponse represents the response from ICX memory recall.
// LocalFallback is true when matches came from the process-local store, not hosted ICX.
type RecallResponse struct {
	Success       bool                  `json:"success"`
	SpaceID       string                `json:"space_id"`
	Query         string                `json:"query"`
	Matches       []MemoryRegisterMatch `json:"matches"`
	LatticeNodes  int                   `json:"lattice_nodes,omitempty"`
	LatencyMs     float64               `json:"latency_ms"`
	Grounding     float64               `json:"grounding_score"`
	LocalFallback bool                  `json:"local_fallback,omitempty"`
	Error         string                `json:"error,omitempty"`
}

// ChatMessage represents a standard chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents a request to the ICX Chat Bridge
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	Stream      bool          `json:"stream"`
	Tools       any           `json:"tools,omitempty"`
}

// ChatCompletionResponse represents a response from the ICX Chat Bridge
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ThoughtSummary   string `json:"thought_summary,omitempty"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	VLNTelemetry struct {
		RecallLatencyMs float64 `json:"recall_latency_ms"`
		FactsRetrieved  int     `json:"facts_retrieved"`
		GroundingScore  float64 `json:"grounding_score"`
		LLMStatus       string  `json:"llm_status"`
	} `json:"vln_telemetry"`
	Error string `json:"error,omitempty"`
}

// TelemetryRecord captures exact benchmark telemetry for an execution
type TelemetryRecord struct {
	Timestamp      time.Time `json:"timestamp"`
	TaskID         string    `json:"task_id"`
	Mode           string    `json:"mode"`
	PromptTokens   int       `json:"prompt_tokens"`
	ResponseTokens int       `json:"response_tokens"`
	LatencyMs      float64   `json:"latency_ms"`
	ToolsProvided  int       `json:"tools_provided"`
	ToolCalled     string    `json:"tool_called,omitempty"`
	Success        bool      `json:"success"`
	IsRefusal      bool      `json:"is_refusal"`
	RefusalCorrect bool      `json:"refusal_correct"`
}
