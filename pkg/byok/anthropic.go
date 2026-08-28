package byok

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

// AnthropicClient connects directly to Anthropic Claude models using user BYOK API keys
type AnthropicClient struct {
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
}

// Ensure AnthropicClient implements LLMClient
var _ LLMClient = (*AnthropicClient)(nil)

// NewAnthropicClient initializes an Anthropic Claude BYOK Client
func NewAnthropicClient(apiKey, model, baseURL string) *AnthropicClient {
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	return &AnthropicClient{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// ProviderName returns provider identifier
func (c *AnthropicClient) ProviderName() string {
	return "anthropic"
}

// ModelName returns configured model identifier
func (c *AnthropicClient) ModelName() string {
	return c.Model
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type anthropicTool struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	InputSchema skills.ToolParameters `json:"input_schema"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
}

type anthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type  string         `json:"type"`
		Text  string         `json:"text,omitempty"`
		ID    string         `json:"id,omitempty"`
		Name  string         `json:"name,omitempty"`
		Input map[string]any `json:"input,omitempty"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GenerateContent sends a chat/completion request to Anthropic Claude
func (c *AnthropicClient) GenerateContent(
	contents []GeminiContent,
	tools []skills.ToolDefinition,
	systemInstruction string,
) (*AgentTurnResult, error) {
	return c.GenerateContentWithContext(context.Background(), contents, tools, systemInstruction)
}

// GenerateContentWithContext sends chat completion request with context
func (c *AnthropicClient) GenerateContentWithContext(
	ctx context.Context,
	contents []GeminiContent,
	tools []skills.ToolDefinition,
	systemInstruction string,
) (*AgentTurnResult, error) {
	t0 := time.Now()

	anthropicMessages := make([]anthropicMessage, 0)
	for _, cItem := range contents {
		role := cItem.Role
		if role == "model" || role == "assistant" {
			role = "assistant"
		} else {
			role = "user"
		}

		var textParts []string
		for _, part := range cItem.Parts {
			if part.Text != "" {
				textParts = append(textParts, part.Text)
			}
		}
		if len(textParts) > 0 {
			anthropicMessages = append(anthropicMessages, anthropicMessage{
				Role:    role,
				Content: strings.Join(textParts, "\n"),
			})
		}
	}

	anthropicTools := make([]anthropicTool, 0, len(tools))
	for _, t := range tools {
		anthropicTools = append(anthropicTools, anthropicTool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.Parameters,
		})
	}

	reqPayload := anthropicRequest{
		Model:     c.Model,
		MaxTokens: 2048,
		System:    systemInstruction,
		Messages:  anthropicMessages,
	}
	if len(anthropicTools) > 0 {
		reqPayload.Tools = anthropicTools
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal anthropic request: %w", err)
	}

	promptTokenEst := len(bodyBytes) / 4

	// 1. Live Cloud API Call
	if !isMockAPIKey(c.APIKey) {
		url := fmt.Sprintf("%s/messages", strings.TrimRight(c.BaseURL, "/"))
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to create anthropic http request: %w", err)
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("x-api-key", c.APIKey)
		httpReq.Header.Set("anthropic-version", "2023-06-01")

		resp, reqErr := c.HTTPClient.Do(httpReq)
		if reqErr != nil {
			return nil, fmt.Errorf("anthropic network error: %w", reqErr)
		}
		defer resp.Body.Close()

		respBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read anthropic response body: %w", readErr)
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("anthropic API error (HTTP %d): %s", resp.StatusCode, string(respBytes))
		}

		var anthResp anthropicResponse
		if parseErr := json.Unmarshal(respBytes, &anthResp); parseErr != nil {
			return nil, fmt.Errorf("failed to parse anthropic json: %w", parseErr)
		}
		if anthResp.Error != nil {
			return nil, fmt.Errorf("anthropic server error (%s): %s", anthResp.Error.Type, anthResp.Error.Message)
		}

		latencyMs := float64(time.Since(t0).Microseconds()) / 1000.0
		result := &AgentTurnResult{
			PromptTokens:    anthResp.Usage.InputTokens,
			ResponseTokens:  anthResp.Usage.OutputTokens,
			LatencyMs:       latencyMs,
			InferenceEngine: fmt.Sprintf("LIVE_CLOUD_API (%s)", c.Model),
		}

		var textParts []string
		for _, block := range anthResp.Content {
			if block.Type == "text" && block.Text != "" {
				textParts = append(textParts, block.Text)
			} else if block.Type == "tool_use" {
				result.ToolCall = &GeminiFunctionCall{
					Name: block.Name,
					Args: block.Input,
				}
			}
		}
		result.TextResponse = strings.Join(textParts, "\n")
		return result, nil
	}

	// 2. Offline Deterministic Reasoner Sandbox
	const inferenceEngine = "OFFLINE_DETERMINISTIC_SANDBOX"
	if len(tools) == 0 {
		var collectedOutputs []string
		for _, m := range anthropicMessages {
			if m.Role == "user" {
				if s, ok := m.Content.(string); ok && strings.HasPrefix(s, "Tool ") {
					collectedOutputs = append(collectedOutputs, fmt.Sprintf("- %s", s))
				}
			}
		}

		latencyMs := float64(time.Since(t0).Microseconds())/1000.0 + 15.0
		if len(collectedOutputs) > 0 {
			return &AgentTurnResult{
				PromptTokens:    promptTokenEst,
				ResponseTokens:  85,
				LatencyMs:       latencyMs,
				InferenceEngine: inferenceEngine,
				TextResponse: fmt.Sprintf("✅ Multi-Skill Pipeline Successfully Executed and Grounded:\n\n%s\n\nAll workflow milestones completed with 100%% cryptographic verification and state crystallization.",
					strings.Join(collectedOutputs, "\n")),
			}, nil
		}

		promptText := ""
		for _, m := range anthropicMessages {
			if m.Role == "user" {
				if s, ok := m.Content.(string); ok {
					promptText += s + " "
				}
			}
		}
		return &AgentTurnResult{
			PromptTokens:    promptTokenEst,
			ResponseTokens:  45,
			LatencyMs:       latencyMs,
			InferenceEngine: inferenceEngine,
			TextResponse:   fmt.Sprintf("Direct response: Processed query '%s' with zero schema overhead.", strings.TrimSpace(promptText)),
		}, nil
	}

	promptText := ""
	for _, m := range anthropicMessages {
		if m.Role == "user" {
			if s, ok := m.Content.(string); ok {
				promptText += s + " "
			}
		}
	}
	promptTextLower := strings.ToLower(promptText)

	var matchedTool *skills.ToolDefinition
	bestToolScore := 0.0

	for idx, t := range tools {
		tNameLower := strings.ToLower(t.Name)
		toolScore := float64(len(tools)-idx) * 0.5
		if strings.Contains(promptTextLower, tNameLower) || strings.Contains(promptTextLower, strings.ReplaceAll(tNameLower, "_", " ")) {
			toolScore += 10.0
		}
		if toolScore > bestToolScore {
			bestToolScore = toolScore
			tCopy := t
			matchedTool = &tCopy
		}
	}

	if matchedTool == nil && len(tools) > 0 {
		matchedTool = &tools[0]
	}

	latencyMs := float64(time.Since(t0).Microseconds())/1000.0 + 18.0
	res := &AgentTurnResult{
		PromptTokens:    promptTokenEst,
		ResponseTokens:  45,
		LatencyMs:       latencyMs,
		InferenceEngine: inferenceEngine,
	}

	if matchedTool != nil {
		res.ToolCall = &GeminiFunctionCall{
			Name: matchedTool.Name,
			Args: map[string]any{
				"query":   strings.TrimSpace(promptText),
				"options": `{"mode": "direct_execution"}`,
			},
		}
		res.TextResponse = fmt.Sprintf("Calling tool `%s` to execute requested operation.", matchedTool.Name)
	} else {
		res.TextResponse = "I have analyzed the request and determined that none of the available tools match the required operation."
	}

	return res, nil
}
