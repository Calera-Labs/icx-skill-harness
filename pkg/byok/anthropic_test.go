package byok

import (
	"context"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func TestAnthropicMultiTurnToolResult(t *testing.T) {
	client := NewAnthropicClient("mock", "claude-3-5-sonnet-20241022", "")

	contents := []GeminiContent{
		{
			Role:  "user",
			Parts: []GeminiPart{{Text: "Extract SEC filing and calculate DCF"}},
		},
		{
			Role: "model",
			Parts: []GeminiPart{
				{
					FunctionCall: &GeminiFunctionCall{
						Name: "sec_edgar_query",
						Args: map[string]any{"ticker": "AAPL"},
					},
				},
			},
		},
		{
			Role: "user",
			Parts: []GeminiPart{
				{
					FunctionResponse: &GeminiFunctionResponse{
						Name:     "sec_edgar_query",
						Response: map[string]any{"operating_margin": 0.3125},
					},
					Text: "Tool sec_edgar_query output: 0.3125",
				},
			},
		},
	}

	// Turn with 0 active tools -> synthesis
	res, err := client.GenerateContentWithContext(context.Background(), contents, nil, "System prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TextResponse == "" {
		t.Errorf("expected non-empty text response for synthesis turn")
	}
}

func TestAnthropicToolConversion(t *testing.T) {
	client := NewAnthropicClient("mock", "claude-3-7-sonnet", "")
	if client.ModelName() != "claude-3-7-sonnet" {
		t.Errorf("expected claude-3-7-sonnet, got %s", client.ModelName())
	}
	tools := []skills.ToolDefinition{
		{
			Name:        "test_tool",
			Description: "A test tool",
			Parameters: skills.ToolParameters{
				Type: "object",
				Properties: map[string]skills.ParameterProperty{
					"query": {Type: "string"},
				},
			},
		},
	}
	contents := []GeminiContent{
		{
			Role:  "user",
			Parts: []GeminiPart{{Text: "Run test tool"}},
		},
	}
	res, err := client.GenerateContentWithContext(context.Background(), contents, tools, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ToolCall == nil || res.ToolCall.Name != "test_tool" {
		t.Errorf("expected test_tool call, got %v", res.ToolCall)
	}
}
