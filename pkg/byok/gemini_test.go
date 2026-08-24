package byok

import (
	"fmt"
	"os"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func TestGeminiDirectCall(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping live Gemini API integration test: GEMINI_API_KEY environment variable not set")
	}

	client := NewGeminiClient(apiKey, "gemini-3.6-flash")

	contents := []GeminiContent{
		{
			Role: "user",
			Parts: []GeminiPart{
				{Text: "Reconcile customer invoice in_99218a with Stripe payment intent and webhook."},
			},
		},
	}

	testTools := []skills.ToolDefinition{
		{
			Name:        "stripe_reconciler",
			Description: "Reconcile Stripe customer invoices and webhooks",
			Category:    "FinTech",
			Parameters: skills.ToolParameters{
				Type: "object",
				Properties: map[string]skills.ParameterProperty{
					"query": {
						Type:        "string",
						Description: "Primary action query",
					},
					"options": {
						Type:        "string",
						Description: "Optional parameters",
					},
				},
				Required: []string{"query"},
			},
		},
	}

	res, err := client.GenerateContent(contents, testTools, "You are a specialized agent. Use the provided tools.")
	if err != nil {
		t.Fatalf("GenerateContent failed: %v", err)
	}

	fmt.Printf("=== GEMINI TEST RESPONSE ===\n")
	fmt.Printf("TextResponse: %s\n", res.TextResponse)
	if res.ToolCall != nil {
		fmt.Printf("ToolCall Name: %s, Args: %+v\n", res.ToolCall.Name, res.ToolCall.Args)
	} else {
		fmt.Printf("ToolCall: nil\n")
	}
	fmt.Printf("PromptTokens: %d, ResponseTokens: %d, LatencyMs: %.2f\n", res.PromptTokens, res.ResponseTokens, res.LatencyMs)
}
