package byok

import (
	"context"
	"strings"

	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

// LLMClient represents a universal BYOK client interface for all model providers
type LLMClient interface {
	GenerateContent(contents []GeminiContent, tools []skills.ToolDefinition, systemInstruction string) (*AgentTurnResult, error)
	GenerateContentWithContext(ctx context.Context, contents []GeminiContent, tools []skills.ToolDefinition, systemInstruction string) (*AgentTurnResult, error)
	ProviderName() string
	ModelName() string
}

// GeminiContent represents content in Gemini API
type GeminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of the content (text or function call/response)
type GeminiPart struct {
	Text             string                  `json:"text,omitempty"`
	FunctionCall     *GeminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *GeminiFunctionResponse `json:"functionResponse,omitempty"`
	Thought          string                  `json:"thought,omitempty"`
}

// GeminiFunctionCall represents a tool call requested by the model
type GeminiFunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}

// GeminiFunctionResponse represents the execution result of a tool call
type GeminiFunctionResponse struct {
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

// GeminiTool wraps function declarations for Gemini
type GeminiTool struct {
	FunctionDeclarations []GeminiFunctionDeclaration `json:"functionDeclarations"`
}

// GeminiFunctionDeclaration represents a callable tool schema for Gemini
type GeminiFunctionDeclaration struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  skills.ToolParameters `json:"parameters"`
}

// GeminiGenerateRequest payload for Google Gemini API
type GeminiGenerateRequest struct {
	Contents          []GeminiContent  `json:"contents"`
	Tools             []GeminiTool     `json:"tools,omitempty"`
	SystemInstruction *GeminiContent   `json:"systemInstruction,omitempty"`
	GenerationConfig  *GeminiGenConfig `json:"generationConfig,omitempty"`
}

// GeminiGenConfig generation parameters
type GeminiGenConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	TopP            float64 `json:"topP,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GeminiGenerateResponse response payload from Google Gemini API
type GeminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []GeminiPart `json:"parts"`
			Role  string       `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

// OpenAIMessage represents an OpenAI format message
type OpenAIMessage struct {
	Role       string            `json:"role"`
	Content    string            `json:"content,omitempty"`
	ToolCalls  []OpenAIToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string            `json:"tool_call_id,omitempty"`
}

// OpenAIToolCall represents a tool call in OpenAI format
type OpenAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// OpenAITool represents a tool schema for OpenAI-compatible APIs
type OpenAITool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string                `json:"name"`
		Description string                `json:"description"`
		Parameters  skills.ToolParameters `json:"parameters"`
	} `json:"function"`
}

// OpenAIChatRequest payload for OpenAI-compatible endpoints
type OpenAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Tools       []OpenAITool    `json:"tools,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

// OpenAIChatResponse payload from OpenAI-compatible endpoints
type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    any    `json:"code"`
	} `json:"error,omitempty"`
}

// AgentTurnResult holds the result of a single agent conversation turn
type AgentTurnResult struct {
	TextResponse    string                `json:"text_response"`
	ToolCall        *GeminiFunctionCall   `json:"tool_call,omitempty"`
	PromptTokens    int                   `json:"prompt_tokens"`
	ResponseTokens  int                   `json:"response_tokens"`
	LatencyMs       float64               `json:"latency_ms"`
	IsRefusal       bool                  `json:"is_refusal"`
	RefusalReason   string                `json:"refusal_reason,omitempty"`
	InferenceEngine string                `json:"inference_engine,omitempty"`
	Viewport        *skills.SkillViewport `json:"viewport,omitempty"`
}

// isMockAPIKey reports whether an API key represents a local test/offline sandbox fixture
func isMockAPIKey(key string) bool {
	k := strings.ToLower(strings.TrimSpace(key))
	return k == "" || k == "mock" || strings.HasPrefix(k, "mock") || strings.Contains(k, "mock") || strings.Contains(k, "offline") || strings.HasPrefix(k, "sk-ant-mock") || strings.HasPrefix(k, "aizasy_mock")
}

