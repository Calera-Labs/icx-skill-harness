package benchmark

import (
	"context"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/agent"
	"github.com/caleralabs/icx-skill-harness/pkg/byok"
	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/router"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func TestBenchmarkSuiteRun(t *testing.T) {
	reg := skills.NewSkillRegistry()

	// Register sample skills for test cases
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
		Instructions: "Use sec_edgar_query",
		MerkleSeal:   "seal_sec",
	})

	icxClient := icx.NewClient(icx.Config{
		SpaceID:       "test_space",
		LocalFallback: true,
	})

	jitRouter := router.NewLatticeSkillRouter(reg, icxClient, router.DefaultRouterConfig())
	llmClient := byok.NewGeminiClient("AIzaSy_mock", "gemini-3.5-flash-lite")
	runner := agent.NewRunner(reg, jitRouter, llmClient, icxClient, "test_space")

	suite := NewSuite(runner, reg)
	cases := suite.GenerateTestCases()

	if len(cases) == 0 {
		t.Fatalf("expected non-empty test cases")
	}

	result, err := suite.RunWithContext(context.Background())
	if err != nil {
		t.Fatalf("unexpected benchmark suite run error: %v", err)
	}

	if result.TotalCases != len(cases) {
		t.Errorf("expected %d total cases, got %d", len(cases), result.TotalCases)
	}
	if result.ICXRefusalAccuracy != 100.0 {
		t.Errorf("expected 100%% ICX refusal accuracy on unanswerable traps, got %.2f%%", result.ICXRefusalAccuracy)
	}
	if result.TokenSavingsPct <= 0 {
		t.Errorf("expected positive token savings percentage")
	}

	jsonBytes, err := result.ExportJSON()
	if err != nil || len(jsonBytes) == 0 {
		t.Errorf("expected successful JSON export of benchmark results")
	}
}
