package gateway

import (
	"sync/atomic"
	"time"
)

// GatewayConfig holds runtime settings for the ICX proxy daemon
type GatewayConfig struct {
	Port                  int           `json:"port"`
	Host                  string        `json:"host"`
	SpaceID               string        `json:"space_id"`
	DefaultProvider       string        `json:"default_provider"`
	DefaultModel          string        `json:"default_model"`
	BYOKKey               string        `json:"-"`
	BYOKBaseURL           string        `json:"byok_base_url"`
	MaxToolsPerViewport   int           `json:"max_tools_per_viewport"`
	OfflineMode           bool          `json:"offline_mode"`
	EnableCrystallization bool          `json:"enable_crystallization"`
	SkillsDir             string        `json:"skills_dir"`
	MCPConfigFile         string        `json:"mcp_config_file"`
	OpenAPIFile           string        `json:"openapi_file"`
	Timeout               time.Duration `json:"timeout"`
	GatewayToken          string        `json:"-"`
	AllowPublicBind       bool          `json:"allow_public_bind"`
	EnableHotRegister     bool          `json:"enable_hot_register"`
	MaxBodyBytes          int64         `json:"max_body_bytes"`
	MaxIncomingTools      int           `json:"max_incoming_tools"`
}

// DefaultGatewayConfig returns standard production defaults
func DefaultGatewayConfig() GatewayConfig {
	return GatewayConfig{
		Port:                  8080,
		Host:                  "127.0.0.1",
		SpaceID:               "local",
		DefaultProvider:       "gemini",
		DefaultModel:          "gemini-3.5-flash-lite",
		MaxToolsPerViewport:   2,
		OfflineMode:           false,
		EnableCrystallization: false,
		Timeout:               30 * time.Second,
		MaxBodyBytes:          defaultMaxBodyBytes,
		MaxIncomingTools:      defaultMaxTools,
	}
}

// OpenAIChatCompletionRequest defines the standard OpenAI chat completions payload
type OpenAIChatCompletionRequest struct {
	Model            string          `json:"model"`
	Messages         []OpenAIMessage `json:"messages"`
	Tools            []OpenAITool    `json:"tools,omitempty"`
	ToolChoice       any             `json:"tool_choice,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	N                int             `json:"n,omitempty"`
	Stream           bool            `json:"stream,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	PresencePenalty  float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64         `json:"frequency_penalty,omitempty"`
	User             string          `json:"user,omitempty"`
	// Calera ICX metadata extensions
	ICXSpaceID      string `json:"icx_space_id,omitempty"`
	ICXSessionID    string `json:"icx_session_id,omitempty"`
	ICXDisablePrune bool   `json:"icx_disable_prune,omitempty"`
}

// OpenAIMessage defines a single message in the conversation
type OpenAIMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content"`
	Name       string           `json:"name,omitempty"`
	ToolCalls  []OpenAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

// OpenAITool defines a callable function tool
type OpenAITool struct {
	Type     string            `json:"type"` // "function"
	Function OpenAIFunctionDef `json:"function"`
}

// OpenAIFunctionDef defines the function signature in OpenAI format
type OpenAIFunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// OpenAIToolCall defines a tool call emitted by the LLM
type OpenAIToolCall struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // "function"
	Function OpenAIFunctionCallData `json:"function"`
}

// OpenAIFunctionCallData holds name and stringified JSON arguments
type OpenAIFunctionCallData struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// OpenAIChatCompletionResponse defines the standard OpenAI response
type OpenAIChatCompletionResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
	Choices           []OpenAIChoice `json:"choices"`
	Usage             OpenAIUsage    `json:"usage"`
}

// OpenAIChoice represents a completion alternative
type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// OpenAIChatCompletionChunk defines a streaming chunk for OpenAI SSE protocol
type OpenAIChatCompletionChunk struct {
	ID                string              `json:"id"`
	Object            string              `json:"object"`
	Created           int64               `json:"created"`
	Model             string              `json:"model"`
	SystemFingerprint string              `json:"system_fingerprint,omitempty"`
	Choices           []OpenAIChunkChoice `json:"choices"`
	Usage             *OpenAIUsage        `json:"usage,omitempty"`
}

// OpenAIChunkChoice defines choices in a streaming chunk
type OpenAIChunkChoice struct {
	Index        int              `json:"index"`
	Delta        OpenAIChunkDelta `json:"delta"`
	FinishReason *string          `json:"finish_reason"`
}

// OpenAIChunkDelta defines the content delta in a streaming chunk
type OpenAIChunkDelta struct {
	Role      string                `json:"role,omitempty"`
	Content   string                `json:"content,omitempty"`
	ToolCalls []OpenAIChunkToolCall `json:"tool_calls,omitempty"`
}

// OpenAIChunkToolCall defines a tool call delta in streaming chunk
type OpenAIChunkToolCall struct {
	Index    int                     `json:"index"`
	ID       string                  `json:"id,omitempty"`
	Type     string                  `json:"type,omitempty"`
	Function OpenAIChunkFunctionData `json:"function"`
}

// OpenAIChunkFunctionData holds streaming tool function arguments
type OpenAIChunkFunctionData struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// OpenAIUsage represents token counts with ICX cost savings metrics
type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
	// ICX Optimization Extensions
	MonolithicTokens int     `json:"icx_monolithic_tokens,omitempty"`
	ICXTokensSaved   int     `json:"icx_tokens_saved,omitempty"`
	ICXSavingsPct    float64 `json:"icx_savings_pct,omitempty"`
	ICXMerkleSeal    string  `json:"icx_merkle_seal,omitempty"`
	ICXSpaceID       string  `json:"icx_space_id,omitempty"`
}

// OpenAIModelListResponse defines standard /v1/models response
type OpenAIModelListResponse struct {
	Object string            `json:"object"`
	Data   []OpenAIModelItem `json:"data"`
}

// OpenAIModelItem defines a single model descriptor
type OpenAIModelItem struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// GatewayStats tracks real-time traffic, token reduction, and dollar savings
type GatewayStats struct {
	TotalRequests     atomic.Int64
	TotalToolsPruned  atomic.Int64
	TotalTokensSaved  atomic.Int64
	TotalTokensServed atomic.Int64
	TotalStatesStored atomic.Int64
	StartTime         time.Time
}

// StatsSnapshot represents a point-in-time snapshot of gateway metrics
type StatsSnapshot struct {
	TotalRequests        int64   `json:"total_requests"`
	TotalToolsPruned     int64   `json:"total_tools_pruned"`
	TotalTokensSaved     int64   `json:"total_tokens_saved"`
	TotalTokensServed    int64   `json:"total_tokens_served"`
	EstimatedDollarSaved float64 `json:"estimated_dollar_saved"`
	TotalStatesStored    int64   `json:"total_states_crystallized"`
	OverallSavingsPct    float64 `json:"overall_savings_pct"`
	UptimeSeconds        int64   `json:"uptime_seconds"`
}
