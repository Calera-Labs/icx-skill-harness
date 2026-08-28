package byok

import (
	"context"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func TestGeminiClientLocalReasoner(t *testing.T) {
	client := NewGeminiClient("AIzaSy_mock_offline_key", "gemini-3.5-flash-lite")

	tools := []skills.ToolDefinition{
		{
			Name:        "stripe_reconciler",
			Description: "Reconcile Stripe customer invoices and payments",
			Parameters: skills.ToolParameters{
				Type: "object",
				Properties: map[string]skills.ParameterProperty{
					"query": {Type: "string"},
				},
				Required: []string{"query"},
			},
			Category: "FinTech",
		},
	}

	contents := []GeminiContent{
		{
			Role:  "user",
			Parts: []GeminiPart{{Text: "Reconcile customer invoice in_99218a with Stripe payment intent"}},
		},
	}

	res, err := client.GenerateContentWithContext(context.Background(), contents, tools, "You are a helpful assistant.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ToolCall == nil {
		t.Fatalf("expected tool call to be generated")
	}
	if res.ToolCall.Name != "stripe_reconciler" {
		t.Errorf("expected 'stripe_reconciler', got '%s'", res.ToolCall.Name)
	}
	if res.PromptTokens <= 0 {
		t.Errorf("expected positive prompt tokens")
	}
}

func TestOpenAIClientLocalReasoner(t *testing.T) {
	client := NewOpenAIClient("mock_openai_key", "gpt-4o-mini", "http://127.0.0.1:59999/unreachable")

	tools := []skills.ToolDefinition{
		{
			Name:        "docker_manager",
			Description: "Inspect and manage Docker container lifecycle",
			Parameters: skills.ToolParameters{
				Type: "object",
				Properties: map[string]skills.ParameterProperty{
					"query": {Type: "string"},
				},
				Required: []string{"query"},
			},
			Category: "DevOps",
		},
	}

	contents := []GeminiContent{
		{
			Role:  "user",
			Parts: []GeminiPart{{Text: "Restart the lithos-trader-01 microservice in docker"}},
		},
	}

	res, err := client.GenerateContentWithContext(context.Background(), contents, tools, "You are a helpful assistant.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ToolCall == nil {
		t.Fatalf("expected tool call to be generated")
	}
	if res.ToolCall.Name != "docker_manager" {
		t.Errorf("expected 'docker_manager', got '%s'", res.ToolCall.Name)
	}
}

func TestAnthropicClientLocalReasoner(t *testing.T) {
	client := NewAnthropicClient("sk-ant-mock", "claude-3-5-sonnet-20241022", "http://127.0.0.1:59999/unreachable")

	tools := []skills.ToolDefinition{
		{
			Name:        "sec_edgar_query",
			Description: "Query SEC EDGAR database for verified financial metrics",
			Parameters: skills.ToolParameters{
				Type: "object",
				Properties: map[string]skills.ParameterProperty{
					"query": {Type: "string"},
				},
				Required: []string{"query"},
			},
			Category: "Finance",
		},
	}

	contents := []GeminiContent{
		{
			Role:  "user",
			Parts: []GeminiPart{{Text: "Extract Apple FY2025 operating margin from SEC 10-K filing"}},
		},
	}

	res, err := client.GenerateContentWithContext(context.Background(), contents, tools, "You are a helpful assistant.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ToolCall == nil {
		t.Fatalf("expected tool call to be generated")
	}
	if res.ToolCall.Name != "sec_edgar_query" {
		t.Errorf("expected 'sec_edgar_query', got '%s'", res.ToolCall.Name)
	}
	if client.ProviderName() != "anthropic" {
		t.Errorf("expected provider anthropic, got %s", client.ProviderName())
	}
}

func TestProviderFactory(t *testing.T) {
	p1 := NewProvider("gemini", "key1", "gemini-3.5-flash-lite", "")
	if p1.ProviderName() != "gemini" {
		t.Errorf("expected provider gemini, got %s", p1.ProviderName())
	}

	p2 := NewProvider("openai", "key2", "gpt-4o", "")
	if p2.ProviderName() != "openai" {
		t.Errorf("expected provider openai, got %s", p2.ProviderName())
	}

	p3 := NewProvider("deepseek", "key3", "deepseek-chat", "")
	if p3.ProviderName() != "openai" { // deepseek uses OpenAI-compatible client
		t.Errorf("expected deepseek to use openai client interface")
	}

	p4 := NewProvider("anthropic", "key4", "claude-3-5-sonnet", "")
	if p4.ProviderName() != "anthropic" {
		t.Errorf("expected provider anthropic, got %s", p4.ProviderName())
	}
}
