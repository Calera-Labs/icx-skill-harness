package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/byok"
	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/router"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func setupTestRunner() (*Runner, *skills.SkillRegistry) {
	reg := skills.NewSkillRegistry()

	reg.Register(&skills.Skill{
		ID:          "sec_edgar_analyst",
		Name:        "SEC EDGAR Financial Analyst",
		Category:    "Finance",
		Description: "Retrieves 10-K, 10-Q filings and exact GAAP metrics from SEC EDGAR",
		Triggers:    []string{"sec", "10-k", "10-q", "gaap", "operating margin"},
		Tools: []skills.ToolDefinition{
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
		},
		Instructions: "Use sec_edgar_query for financial data.",
		MerkleSeal:   "seal_sec",
	})

	icxClient := icx.NewClient(icx.Config{
		SpaceID:       "test_space",
		LocalFallback: true,
	})

	cfg := router.DefaultRouterConfig()
	jitRouter := router.NewLatticeSkillRouter(reg, icxClient, cfg)
	llmClient := byok.NewGeminiClient("AIzaSy_mock", "gemini-3.5-flash-lite")

	runner := NewRunner(reg, jitRouter, llmClient, icxClient, "test_space")
	return runner, reg
}

func TestRunnerExecuteWithICX(t *testing.T) {
	runner, _ := setupTestRunner()

	res, err := runner.ExecuteWithICX("Fetch Apple operating margin from SEC 10-K", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.IsRefusal {
		t.Fatalf("expected successful execution, got refusal")
	}
	if res.ToolCall == nil || res.ToolCall.Name != "sec_edgar_query" {
		t.Errorf("expected tool call 'sec_edgar_query'")
	}
}

func TestRunnerRefusalOnTrap(t *testing.T) {
	runner, _ := setupTestRunner()

	res, err := runner.ExecuteWithICX("Compile Rust smart contract on Solana mainnet", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.IsRefusal {
		t.Errorf("expected safe refusal on unanswerable trap")
	}
	if !strings.Contains(res.TextResponse, "SAFE_REFUSAL") {
		t.Errorf("expected safe refusal message in text response")
	}
}

func TestRunnerAgentLoop(t *testing.T) {
	runner, _ := setupTestRunner()

	trace, err := runner.RunAgentLoop(context.Background(), "Fetch Apple operating margin from SEC 10-K", "", 3, nil)
	if err != nil {
		t.Fatalf("unexpected error running agent loop: %v", err)
	}

	if trace.TurnsExecuted < 1 {
		t.Errorf("expected at least 1 turn executed")
	}
	if trace.Turns[0].ToolCall == nil {
		t.Errorf("expected turn 1 to execute tool call")
	}
	if trace.Turns[0].ToolOutput == "" {
		t.Errorf("expected tool output to be recorded")
	}
}

func TestCrystallizeToolOutput(t *testing.T) {
	runner, _ := setupTestRunner()

	// Short output is kept inline
	shortOut, _ := runner.CrystallizeToolOutput("tool_a", "short result")
	if shortOut != "short result" {
		t.Errorf("expected short output to remain inline")
	}

	// Large output is crystallized
	largeData := strings.Repeat("ABCDEF1234567890\n", 50)
	crysOut, err := runner.CrystallizeToolOutput("sec_edgar", largeData)
	if err != nil {
		t.Fatalf("unexpected error crystallizing: %v", err)
	}

	if !strings.Contains(crysOut, "[LOCAL_FALLBACK_STORED]") {
		t.Errorf("expected local-fallback store reference in output, got: %s", crysOut)
	}
	if strings.Contains(crysOut, "crystallized in ICX") || strings.Contains(strings.ToLower(crysOut), "a4 lattice") {
		t.Errorf("local fallback must not claim a hosted ICX write, got: %s", crysOut)
	}
}
