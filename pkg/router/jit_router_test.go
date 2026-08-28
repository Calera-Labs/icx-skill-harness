package router

import (
	"context"
	"strings"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func setupTestRegistryAndRouter() (*skills.SkillRegistry, *LatticeSkillRouter) {
	reg := skills.NewSkillRegistry()

	// Register sample skills
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
						"ticker": {Type: "string"},
					},
					Required: []string{"ticker"},
				},
				Category: "Finance",
			},
		},
		Instructions: "Use tool sec_edgar_query when requested.",
		MerkleSeal:   "seal_sec",
	})

	reg.Register(&skills.Skill{
		ID:          "git_code_patcher",
		Name:        "Git Code Patcher",
		Category:    "DevOps",
		Description: "Inspects repository AST and generates unified diff git patches",
		Triggers:    []string{"git", "diff", "patch", "refactor"},
		Tools: []skills.ToolDefinition{
			{
				Name:        "git_diff_patcher",
				Description: "Synthesize git diff patch for code refactoring",
				Parameters: skills.ToolParameters{
					Type: "object",
					Properties: map[string]skills.ParameterProperty{
						"file": {Type: "string"},
					},
					Required: []string{"file"},
				},
				Category: "DevOps",
			},
		},
		Instructions: "Use tool git_diff_patcher when requested.",
		MerkleSeal:   "seal_git",
	})

	icxClient := icx.NewClient(icx.Config{
		SpaceID:       "test_space",
		LocalFallback: true,
	})

	cfg := DefaultRouterConfig()
	router := NewLatticeSkillRouter(reg, icxClient, cfg)
	return reg, router
}

func TestRouterExactMatch(t *testing.T) {
	_, router := setupTestRegistryAndRouter()

	viewport, err := router.RouteQueryWithContext(context.Background(), "Fetch Apple operating margin from SEC 10-K")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if viewport.IsRefusal {
		t.Fatalf("expected successful routing, got refusal: %s", viewport.RefusalReason)
	}
	if len(viewport.ActiveTools) == 0 {
		t.Fatalf("expected active tools in viewport")
	}
	if viewport.ActiveTools[0].Name != "sec_edgar_query" {
		t.Errorf("expected 'sec_edgar_query', got '%s'", viewport.ActiveTools[0].Name)
	}
	if viewport.TotalSchemaTokens <= 0 {
		t.Errorf("expected positive schema tokens count")
	}
}

func TestRouterSafeRefusalOnUnanswerableTrap(t *testing.T) {
	_, router := setupTestRegistryAndRouter()

	// Prompt has no matching skill in registry
	viewport, err := router.RouteQueryWithContext(context.Background(), "Execute quantum circuit annealing on IBM Quantum QPU")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !viewport.IsRefusal {
		t.Errorf("expected SAFE_REFUSAL for quantum prompt, but viewport accepted it with %d tools", len(viewport.ActiveTools))
	}
}

func TestRouterEmptyQuery(t *testing.T) {
	_, router := setupTestRegistryAndRouter()

	viewport, err := router.RouteQueryWithContext(context.Background(), "   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !viewport.IsRefusal {
		t.Errorf("expected refusal on empty query")
	}
}

func TestRouterFutureSECFilingRefusal(t *testing.T) {
	_, router := setupTestRegistryAndRouter()

	viewport, err := router.RouteQueryWithContext(context.Background(), "Fetch Apple's FY2035 SEC 10-K filing operating margins and net profit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !viewport.IsRefusal {
		t.Errorf("expected SAFE_REFUSAL for future FY2035 SEC filing, but got %d tools", len(viewport.ActiveTools))
	}
}

func TestContainsWordFastMatching(t *testing.T) {
	tests := []struct {
		text   string
		target string
		want   bool
	}{
		{"Fetch SEC filing", "sec", true},
		{"security audit scan", "sec", false},
		{"sec_edgar_query", "sec", true},
		{"check 10-k now", "10-k", true},
		{"unrelated prompt", "git", false},
	}

	for _, tt := range tests {
		got := containsWord(strings.ToLower(tt.text), strings.ToLower(tt.target))
		if got != tt.want {
			t.Errorf("containsWord(%q, %q) = %v, want %v", tt.text, tt.target, got, tt.want)
		}
	}
}

func TestRoutePipelineTurn(t *testing.T) {
	_, router := setupTestRegistryAndRouter()

	prompt := "1. Fetch Apple's FY2025 operating margin from SEC 10-K. 2. Generate a unified git diff patch for src/pricing.py."

	// Turn 1: No tools executed yet -> should route SEC tool
	vp1, err := router.RoutePipelineTurn(context.Background(), prompt, nil)
	if err != nil {
		t.Fatalf("turn 1 error: %v", err)
	}
	if len(vp1.ActiveTools) == 0 || vp1.ActiveTools[0].Name != "sec_edgar_query" {
		t.Errorf("turn 1 expected sec_edgar_query, got %v", vp1.ActiveTools)
	}

	// Turn 2: sec_edgar_query already executed -> should route git tool
	vp2, err := router.RoutePipelineTurn(context.Background(), prompt, []string{"sec_edgar_query"})
	if err != nil {
		t.Fatalf("turn 2 error: %v", err)
	}
	if len(vp2.ActiveTools) == 0 || vp2.ActiveTools[0].Name != "git_diff_patcher" {
		t.Errorf("turn 2 expected git_diff_patcher, got %v", vp2.ActiveTools)
	}

	// Turn 3: both executed -> should return empty active tools for final synthesis
	vp3, err := router.RoutePipelineTurn(context.Background(), prompt, []string{"sec_edgar_query", "git_diff_patcher"})
	if err != nil {
		t.Fatalf("turn 3 error: %v", err)
	}
	if len(vp3.ActiveTools) != 0 {
		t.Errorf("turn 3 expected 0 active tools for final synthesis, got %d", len(vp3.ActiveTools))
	}
}

