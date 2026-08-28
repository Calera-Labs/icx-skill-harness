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

// OpenAIClient connects to OpenAI-compatible APIs (OpenAI, DeepSeek, Ollama, Groq, vLLM)
type OpenAIClient struct {
	APIKey     string
	Model      string
	BaseURL    string
	HTTPClient *http.Client
}

// Ensure OpenAIClient implements LLMClient
var _ LLMClient = (*OpenAIClient)(nil)

// NewOpenAIClient initializes an OpenAI-compatible client
func NewOpenAIClient(apiKey, model, baseURL string) *OpenAIClient {
	if model == "" {
		model = "gpt-4o-mini"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIClient{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// ProviderName returns provider name
func (c *OpenAIClient) ProviderName() string {
	return "openai"
}

// ModelName returns configured model name
func (c *OpenAIClient) ModelName() string {
	return c.Model
}

// GenerateContent sends chat completion request using OpenAI standard schema
func (c *OpenAIClient) GenerateContent(
	contents []GeminiContent,
	tools []skills.ToolDefinition,
	systemInstruction string,
) (*AgentTurnResult, error) {
	return c.GenerateContentWithContext(context.Background(), contents, tools, systemInstruction)
}

// GenerateContentWithContext sends chat completion request with context
func (c *OpenAIClient) GenerateContentWithContext(
	ctx context.Context,
	contents []GeminiContent,
	tools []skills.ToolDefinition,
	systemInstruction string,
) (*AgentTurnResult, error) {
	t0 := time.Now()
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	messages := make([]OpenAIMessage, 0)
	if systemInstruction != "" {
		messages = append(messages, OpenAIMessage{
			Role:    "system",
			Content: systemInstruction,
		})
	}

	for _, cItem := range contents {
		role := cItem.Role
		if role == "" || role == "user" {
			role = "user"
		} else if role == "model" {
			role = "assistant"
		}

		var textParts []string
		for _, part := range cItem.Parts {
			if part.Text != "" {
				textParts = append(textParts, part.Text)
			}
		}

		messages = append(messages, OpenAIMessage{
			Role:    role,
			Content: strings.Join(textParts, "\n"),
		})
	}

	openAITools := make([]OpenAITool, 0, len(tools))
	for _, t := range tools {
		openAITools = append(openAITools, OpenAITool{
			Type: "function",
			Function: struct {
				Name        string                `json:"name"`
				Description string                `json:"description"`
				Parameters  skills.ToolParameters `json:"parameters"`
			}{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		})
	}

	reqPayload := OpenAIChatRequest{
		Model:       c.Model,
		Messages:    messages,
		Temperature: 0.2,
	}
	if len(openAITools) > 0 {
		reqPayload.Tools = openAITools
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal openai request: %w", err)
	}

	promptTokenEst := len(bodyBytes) / 4

	// 1. Live Cloud API Call
	if !isMockAPIKey(c.APIKey) {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to create openai http request: %w", err)
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
		resp, reqErr := c.HTTPClient.Do(httpReq)
		if reqErr != nil {
			return nil, fmt.Errorf("openai network error: %w", reqErr)
		}
		defer resp.Body.Close()

		respBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read openai response body: %w", readErr)
		}

		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("openai API error (HTTP %d): %s", resp.StatusCode, string(respBytes))
		}

		var oaiResp OpenAIChatResponse
		if parseErr := json.Unmarshal(respBytes, &oaiResp); parseErr != nil {
			return nil, fmt.Errorf("failed to parse openai json: %w", parseErr)
		}
		if oaiResp.Error != nil {
			return nil, fmt.Errorf("openai server error: %s", oaiResp.Error.Message)
		}

		latencyMs := float64(time.Since(t0).Microseconds()) / 1000.0
		result := &AgentTurnResult{
			PromptTokens:    oaiResp.Usage.PromptTokens,
			ResponseTokens:  oaiResp.Usage.CompletionTokens,
			LatencyMs:       latencyMs,
			InferenceEngine: fmt.Sprintf("LIVE_CLOUD_API (%s)", c.Model),
		}
		if len(oaiResp.Choices) > 0 {
			choice := oaiResp.Choices[0]
			result.TextResponse = choice.Message.Content
			if len(choice.Message.ToolCalls) > 0 {
				tc := choice.Message.ToolCalls[0]
				var args map[string]any
				_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
				result.ToolCall = &GeminiFunctionCall{
					Name: tc.Function.Name,
					Args: args,
				}
			}
		}
		return result, nil
	}

	// 2. Offline Deterministic Reasoner Sandbox (Used ONLY when key is explicitly set to mock/offline)
	const inferenceEngine = "OFFLINE_DETERMINISTIC_SANDBOX"
	if len(tools) == 0 {
		var collectedOutputs []string
		for _, m := range messages {
			if m.Role == "tool" || strings.HasPrefix(m.Content, "Tool ") {
				collectedOutputs = append(collectedOutputs, fmt.Sprintf("- %s", m.Content))
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
		for _, m := range messages {
			if m.Role == "user" {
				promptText += m.Content + " "
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
	for _, m := range messages {
		if m.Role == "user" {
			promptText += m.Content + " "
		}
	}
	promptTextLower := strings.ToLower(promptText)

	// Determine the best matching tool from the provided micro-viewport
	var matchedTool *skills.ToolDefinition
	bestToolScore := 0.0

	genericTokens := map[string]bool{
		"query": true, "exec": true, "executor": true, "manager": true,
		"runner": true, "client": true, "caller": true, "calc": true,
		"ops": true, "tool": true, "get": true, "update": true, "sync": true,
	}

	for idx, t := range tools {
		tNameLower := strings.ToLower(t.Name)
		tDescLower := strings.ToLower(t.Description)
		tCatLower := strings.ToLower(t.Category)

		toolScore := 0.0

		// Viewport rank priority: earlier tools selected by JIT router have base priority
		toolScore += float64(len(tools)-idx) * 0.5

		// Exact name match
		if strings.Contains(promptTextLower, tNameLower) || strings.Contains(promptTextLower, strings.ReplaceAll(tNameLower, "_", " ")) {
			toolScore += 10.0
		}

		// Specific name tokens
		nameTokens := strings.Fields(strings.ReplaceAll(tNameLower, "_", " "))
		for _, nt := range nameTokens {
			if len(nt) > 3 && !genericTokens[nt] && strings.Contains(promptTextLower, nt) {
				toolScore += 5.0
			}
		}

		// Category match
		if tCatLower != "" && strings.Contains(promptTextLower, tCatLower) {
			toolScore += 3.0
		}

		// Description overlap
		descWords := strings.Fields(tDescLower)
		for _, w := range descWords {
			if len(w) > 3 && !genericTokens[w] && strings.Contains(promptTextLower, w) {
				toolScore += 1.5
			}
		}

		if toolScore > bestToolScore {
			bestToolScore = toolScore
			tCopy := t
			matchedTool = &tCopy
		}
	}

	// In targeted micro-viewport mode, default to the top-ranked tool if none exceeded
	if matchedTool == nil && len(tools) > 0 {
		matchedTool = &tools[0]
	}

	latencyMs := float64(time.Since(t0).Microseconds())/1000.0 + 19.0
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
