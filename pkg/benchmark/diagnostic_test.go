package benchmark

import (
	"context"
	"fmt"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/agent"
	"github.com/caleralabs/icx-skill-harness/pkg/byok"
	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/router"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func setupTestDiagnosticSuite() (*DiagnosticSuite, *skills.SkillRegistry) {
	reg := skills.NewSkillRegistry()
	_ = skills.PopulateCatalog(reg, "")

	icxClient := icx.NewClient(icx.Config{
		SpaceID:       "test_space",
		LocalFallback: true,
	})

	jitRouter := router.NewLatticeSkillRouter(reg, icxClient, router.DefaultRouterConfig())
	llmClient := byok.NewGeminiClient("AIzaSy_mock", "gemini-3.5-flash-lite")
	runner := agent.NewRunner(reg, jitRouter, llmClient, icxClient, "test_space")

	return NewDiagnosticSuite(runner, reg), reg
}

func TestDiagnosticSuiteAllCategories(t *testing.T) {
	suite, _ := setupTestDiagnosticSuite()

	res, err := suite.Run(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected diagnostic suite run error: %v", err)
	}

	for _, cr := range res.CaseResults {
		status := "PASS"
		if !cr.Passed {
			status = fmt.Sprintf("FAIL (%s)", cr.DiagnosedIssue)
		}
		t.Logf("[%s] %-24s -> %s | Tool: '%s' | Lat: %.1fms | %s",
			cr.TestCase.Category, cr.TestCase.ID, status, cr.ToolCalled, cr.LatencyMs, cr.Explanation)
	}

	if res.TotalCases == 0 {
		t.Fatalf("expected non-zero diagnostic test cases")
	}

	if res.OverallPassRate < 80.0 {
		t.Errorf("expected high overall pass rate, got %.2f%%", res.OverallPassRate)
	}

	if res.TokenSavingsPct <= 0.0 {
		t.Errorf("expected positive token savings, got %.2f%%", res.TokenSavingsPct)
	}

	if len(res.ScaleLadderResults) == 0 {
		t.Errorf("expected scale ladder results to be populated")
	}

	report := res.ExportMarkdownReport()
	if report == "" {
		t.Errorf("expected non-empty markdown report")
	}

	jsonData, err := res.ExportJSON()
	if err != nil || len(jsonData) == 0 {
		t.Errorf("expected valid JSON export, err: %v", err)
	}
}

func TestDiagnosticSuiteDistractorFilter(t *testing.T) {
	suite, _ := setupTestDiagnosticSuite()

	res, err := suite.Run(context.Background(), CatDistractorCollision)
	if err != nil {
		t.Fatalf("error running distractor category: %v", err)
	}

	for _, cr := range res.CaseResults {
		if cr.TestCase.Category != CatDistractorCollision {
			t.Errorf("expected only distractor collision cases, got %s", cr.TestCase.Category)
		}
	}
}

func TestDiagnosticSuiteFaultInjection(t *testing.T) {
	suite, _ := setupTestDiagnosticSuite()

	res, err := suite.Run(context.Background(), CatFaultInjection)
	if err != nil {
		t.Fatalf("error running fault injection category: %v", err)
	}

	if res.FaultRecoveryRate < 100.0 {
		t.Errorf("expected 100%% fault recovery rate, got %.2f%%", res.FaultRecoveryRate)
	}
}
