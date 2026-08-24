package router

import (
	"context"
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
