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

func setupE2ETestRunner() (*agent.Runner, *skills.SkillRegistry) {
	reg := skills.NewSkillRegistry()
	_ = skills.PopulateCatalog(reg, "")

	icxClient := icx.NewClient(icx.Config{
		SpaceID:       "test_space",
		LocalFallback: true,
	})

	cfg := router.DefaultRouterConfig()
	cfg.EnableLatticeWalk = false
	jitRouter := router.NewLatticeSkillRouter(reg, icxClient, cfg)
	llmClient := byok.NewGeminiClient("AIzaSy_mock", "gemini-3.5-flash-lite")
	runner := agent.NewRunner(reg, jitRouter, llmClient, icxClient, "test_space")

	return runner, reg
}

func TestE2EBenchmarkSuiteRun(t *testing.T) {
	runner, reg := setupE2ETestRunner()
	suite := NewE2ESuite(runner, reg)
	cases := suite.GeneratePipelineTestCases()

	if len(cases) == 0 {
		t.Fatalf("expected non-empty pipeline test cases")
	}

	result, err := suite.RunWithContext(context.Background())
	if err != nil {
		t.Fatalf("unexpected e2e benchmark error: %v", err)
	}

	if result.TotalPipelines != len(cases) {
		t.Errorf("expected %d pipelines, got %d", len(cases), result.TotalPipelines)
	}
	if result.ICXPipelinePassRate < 90.0 {
		t.Errorf("expected >=90%% ICX pipeline pass rate, got %.2f%%", result.ICXPipelinePassRate)
	}
	if result.MidStreamRefusalAccuracy != 100.0 {
		t.Errorf("expected 100%% mid-stream refusal accuracy on traps, got %.2f%%", result.MidStreamRefusalAccuracy)
	}
	if result.TokenSavingsPct < 85.0 {
		t.Errorf("expected >=85%% token savings, got %.2f%%", result.TokenSavingsPct)
	}

	jsonBytes, err := result.ExportJSON()
	if err != nil || len(jsonBytes) == 0 {
		t.Errorf("expected valid JSON export of E2E benchmark results")
	}
}

func TestE2ESinglePipelineFlow(t *testing.T) {
	runner, _ := setupE2ETestRunner()

	prompt := "1. Extract Apple's FY2025 operating margin and GAAP revenue from SEC 10-K filing. " +
		"2. Calculate a DCF financial valuation model with 8.5% WACC and terminal value. " +
		"3. Update PostgreSQL financial settlement ledger to COMMITTED with transaction hash. " +
		"4. Send an executive briefing summary via Slack webhook."

	trace, err := runner.RunAgentLoop(context.Background(), prompt, "", 6, nil)
	if err != nil {
		t.Fatalf("unexpected agent loop error: %v", err)
	}

	if trace.TurnsExecuted < 4 {
		t.Errorf("expected at least 4 turns, got %d", trace.TurnsExecuted)
	}
	if trace.TotalPromptTokens > 6000 {
		t.Errorf("expected lean prompt token consumption (<6000), got %d", trace.TotalPromptTokens)
	}
}
