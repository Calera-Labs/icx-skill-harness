package router

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

// RouterConfig defines configuration for the JIT Skill Router
type RouterConfig struct {
	MaxToolsPerViewport int     `json:"max_tools_per_viewport"`
	MaxTokensThreshold  int     `json:"max_tokens_threshold"`
	MinScoreThreshold   float64 `json:"min_score_threshold"`
	RefusalThreshold    float64 `json:"refusal_threshold"`
	EnableLatticeWalk   bool    `json:"enable_lattice_walk"`
	LatticeWalkTopK     int     `json:"lattice_walk_top_k"`
}

// DefaultRouterConfig returns standard production defaults
func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		MaxToolsPerViewport: 2,
		MaxTokensThreshold:  600,
		MinScoreThreshold:   1.2,
		RefusalThreshold:    1.8,
		EnableLatticeWalk:   true,
		LatticeWalkTopK:     5,
	}
}

// ScoredSkill holds ranking information for a candidate skill
type ScoredSkill struct {
	Skill       *skills.Skill
	Score       float64
	ExactMatch  bool
	MatchTags   []string
	LatticeWalk bool
}

// LatticeSkillRouter provides sub-millisecond JIT micro-viewport routing over hundreds of skills
type LatticeSkillRouter struct {
	registry *skills.SkillRegistry
	icx      *icx.Client
	config   RouterConfig
}

// NewLatticeSkillRouter initializes the router with registry and ICX client
func NewLatticeSkillRouter(reg *skills.SkillRegistry, icxClient *icx.Client, cfg RouterConfig) *LatticeSkillRouter {
	if cfg.MaxToolsPerViewport <= 0 {
		cfg.MaxToolsPerViewport = 2
	}
	if cfg.MaxTokensThreshold <= 0 {
		cfg.MaxTokensThreshold = 600
	}
	if cfg.LatticeWalkTopK <= 0 {
		cfg.LatticeWalkTopK = 5
	}
	return &LatticeSkillRouter{
		registry: reg,
		icx:      icxClient,
		config:   cfg,
	}
}

// RouteQuery analyzes a user prompt and generates a compact micro-viewport
func (r *LatticeSkillRouter) RouteQuery(prompt string) (*skills.SkillViewport, error) {
	return r.RouteQueryWithContext(context.Background(), prompt)
}

// RouteQueryWithContext analyzes prompt with context cancellation
func (r *LatticeSkillRouter) RouteQueryWithContext(ctx context.Context, prompt string) (*skills.SkillViewport, error) {
	t0 := time.Now()
	allSkills := r.registry.GetAllSkills()

	if len(allSkills) == 0 {
		return &skills.SkillViewport{
			IsRefusal:        true,
			RefusalReason:    "404_NO_SKILLS_REGISTERED",
			RoutingLatencyUs: time.Since(t0).Microseconds(),
		}, nil
	}

	trimmedPrompt := strings.TrimSpace(prompt)
	if trimmedPrompt == "" {
		return &skills.SkillViewport{
			IsRefusal:        true,
			RefusalReason:    "EMPTY_QUERY_INTENT",
			RoutingLatencyUs: time.Since(t0).Microseconds(),
		}, nil
	}

	// 1. Tokenize query
	tokens := r.tokenize(trimmedPrompt)

	// 2. Perform ICX Volumetric Lattice Walk if enabled
	latticeBoostMap := make(map[string]float64)
	latticeWalkUsed := false

	if r.config.EnableLatticeWalk && r.icx != nil {
		recallResp, err := r.icx.RecallWithContext(ctx, icx.RecallRequest{
			Query:  trimmedPrompt,
			TopK:   r.config.LatticeWalkTopK,
			Family: "skill.definition",
		})
		if err == nil && recallResp != nil && len(recallResp.Matches) > 0 {
			latticeWalkUsed = true
			for _, match := range recallResp.Matches {
				docTitle := strings.ToLower(match.DocumentTitle)
				docTitle = strings.TrimSuffix(docTitle, ".md")
				docTitle = strings.TrimSuffix(docTitle, ".json")
				if docTitle != "" {
					latticeBoostMap[docTitle] = match.Confidence * 2.5
				}
			}
		}
	}

	// 3. Score skills using BM25 + Tag Matching + Trigram + Lattice Boost
	scored := make([]ScoredSkill, 0, len(allSkills))
	for _, s := range allSkills {
		score, exact, tags := r.scoreSkill(s, tokens, trimmedPrompt)

		// Apply ICX Lattice walk boost if match found
		if boost, hasBoost := latticeBoostMap[strings.ToLower(s.ID)]; hasBoost {
			score += boost
			tags = append(tags, fmt.Sprintf("icx_lattice:%.2f", boost))
		}

		if score >= r.config.MinScoreThreshold {
			scored = append(scored, ScoredSkill{
				Skill:       s,
				Score:       score,
				ExactMatch:  exact,
				MatchTags:   tags,
				LatticeWalk: latticeWalkUsed,
			})
		}
	}

	// 4. Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].ExactMatch != scored[j].ExactMatch {
			return scored[i].ExactMatch
		}
		return scored[i].Score > scored[j].Score
	})

	// 5. Check for Epistemic Refusal (Missing premise / unanswerable intent)
	if len(scored) == 0 || (len(scored) > 0 && scored[0].Score < r.config.RefusalThreshold) {
		return &skills.SkillViewport{
			IsRefusal:        true,
			RefusalReason:    fmt.Sprintf("SAFE_REFUSAL: No verified skill in lattice matches intent for '%s'", prompt),
			RoutingLatencyUs: time.Since(t0).Microseconds(),
			LatticeWalkUsed:  latticeWalkUsed,
		}, nil
	}

	// 6. Select Top-K skills up to MaxToolsPerViewport and token budget
	selectedSkills := make([]*skills.Skill, 0)
	activeTools := make([]skills.ToolDefinition, 0)
	totalTokens := 0

	for _, sc := range scored {
		if len(activeTools) >= r.config.MaxToolsPerViewport {
			break
		}

		skillToks := sc.Skill.EstimatedTokenSize()
		if totalTokens+skillToks > r.config.MaxTokensThreshold && len(activeTools) > 0 {
			break
		}

		selectedSkills = append(selectedSkills, sc.Skill)
		activeTools = append(activeTools, sc.Skill.Tools...)
		totalTokens += skillToks
	}

	// Measure precise schema token payload
	schemaBytes, _ := json.Marshal(activeTools)
	exactSchemaTokens := len(schemaBytes) / 4

	topScore := 0.0
	if len(scored) > 0 {
		topScore = scored[0].Score
	}

	return &skills.SkillViewport{
		MatchedSkills:     selectedSkills,
		ActiveTools:       activeTools,
		TotalSchemaTokens: exactSchemaTokens,
		ConfidenceScore:   math.Min(1.0, topScore/4.0),
		QueryIntent:       prompt,
		IsRefusal:         false,
		RoutingLatencyUs:  time.Since(t0).Microseconds(),
		LatticeWalkUsed:   latticeWalkUsed,
	}, nil
}

// RoutePipelineTurn routes multi-turn pipeline steps dynamically, updating active micro-viewports per turn
func (r *LatticeSkillRouter) RoutePipelineTurn(
	ctx context.Context,
	initialPrompt string,
	executedTools []string,
) (*skills.SkillViewport, error) {
	t0 := time.Now()
	trimmedPrompt := strings.TrimSpace(initialPrompt)
	if trimmedPrompt == "" {
		return &skills.SkillViewport{
			IsRefusal:        true,
			RefusalReason:    "EMPTY_QUERY_INTENT",
			RoutingLatencyUs: time.Since(t0).Microseconds(),
		}, nil
	}

	// 1. If no tools have been executed yet and not multi-step, route single query
	segments := r.extractStepSegments(trimmedPrompt)
	if len(segments) == 0 {
		if len(executedTools) > 0 {
			// Single-step task with tool already executed -> return empty active tools for final synthesis
			return &skills.SkillViewport{
				MatchedSkills:     nil,
				ActiveTools:       nil,
				TotalSchemaTokens: 0,
				ConfidenceScore:   1.0,
				QueryIntent:       trimmedPrompt,
				IsRefusal:         false,
				RoutingLatencyUs:  time.Since(t0).Microseconds(),
			}, nil
		}
		return r.RouteQueryWithContext(ctx, initialPrompt)
	}

	// 2. Check each step in sequence against executedTools
	executedSet := make(map[string]bool)
	for _, et := range executedTools {
		executedSet[strings.ToLower(et)] = true
	}

	for i, seg := range segments {
		// Route this segment against the lattice
		segViewport, err := r.RouteQueryWithContext(ctx, seg)
		if err != nil {
			return nil, err
		}

		// If a step segment is an unanswerable trap:
		if segViewport.IsRefusal {
			if i <= len(executedTools) {
				return segViewport, nil
			}
		}

		// Check if any tool in this segment has already been executed
		segAlreadyExecuted := false
		for _, tool := range segViewport.ActiveTools {
			if executedSet[strings.ToLower(tool.Name)] {
				segAlreadyExecuted = true
				break
			}
		}

		if !segAlreadyExecuted {
			// Found the next unexecuted pending step
			return segViewport, nil
		}
	}

	// 3. All pipeline segments have been executed: return 0 active tools for final synthesis
	return &skills.SkillViewport{
		MatchedSkills:     nil,
		ActiveTools:       nil,
		TotalSchemaTokens: 0,
		ConfidenceScore:   1.0,
		QueryIntent:       trimmedPrompt,
		IsRefusal:         false,
		RoutingLatencyUs:  time.Since(t0).Microseconds(),
	}, nil
}

// extractStepSegments splits complex prompts into sequential pipeline sub-steps
func (r *LatticeSkillRouter) extractStepSegments(prompt string) []string {
	// 1. Try splitting by numbered list (e.g. "1. ... 2. ... 3. ...")
	numRegex := regexp.MustCompile(`(?i)(?:^|\s)(?:\d+[\.\)]|step\s*\d+[:\.]?)\s+`)
	indices := numRegex.FindAllStringIndex(prompt, -1)
	if len(indices) >= 2 {
		segments := make([]string, 0, len(indices))
		for i := 0; i < len(indices); i++ {
			start := indices[i][1]
			end := len(prompt)
			if i+1 < len(indices) {
				end = indices[i+1][0]
			}
			seg := strings.TrimSpace(prompt[start:end])
			if seg != "" {
				segments = append(segments, seg)
			}
		}
		if len(segments) >= 2 {
			return segments
		}
	}

	// 2. Try splitting by sequence transition connectors ("then", "after that", ";")
	connectorRegex := regexp.MustCompile(`(?i)(?:;\s*|\s+(?:then|after that|afterwards|next|finally|subsequently|and then)\s+)`)
	splitParts := connectorRegex.Split(prompt, -1)
	if len(splitParts) >= 2 {
		segments := make([]string, 0, len(splitParts))
		for _, p := range splitParts {
			pTrim := strings.TrimSpace(p)
			if len(pTrim) > 10 {
				segments = append(segments, pTrim)
			}
		}
		if len(segments) >= 2 {
			return segments
		}
	}

	return nil
}

func (r *LatticeSkillRouter) tokenize(text string) []string {
	cleaned := strings.ToLower(text)
	cleaned = regexp.MustCompile(`[^a-z0-9_\-\.]+`).ReplaceAllString(cleaned, " ")
	words := strings.Fields(cleaned)
	stopwords := map[string]bool{
		"the": true, "and": true, "for": true, "with": true, "this": true,
		"that": true, "from": true, "into": true, "have": true, "what": true,
		"where": true, "when": true, "how": true, "please": true, "want": true,
		"your": true, "can": true, "you": true, "about": true, "more": true,
	}

	res := make([]string, 0, len(words))
	for _, w := range words {
		if len(w) > 2 && !stopwords[w] {
			res = append(res, w)
		}
	}
	return res
}

func (r *LatticeSkillRouter) scoreSkill(skill *skills.Skill, queryTokens []string, rawPrompt string) (float64, bool, []string) {
	score := 0.0
	exact := false
	matchedTags := make([]string, 0)
	rawPromptLower := strings.ToLower(rawPrompt)

	// Check exact name match with boundary
	skillNameLower := strings.ToLower(skill.Name)
	if matchWordBoundary(rawPromptLower, skillNameLower) || matchWordBoundary(rawPromptLower, strings.ReplaceAll(skillNameLower, "_", " ")) {
		score += 5.0
		exact = true
		matchedTags = append(matchedTags, "exact_name")
	}

	// Check skill ID match
	skillIDLower := strings.ToLower(skill.ID)
	if matchWordBoundary(rawPromptLower, skillIDLower) || matchWordBoundary(rawPromptLower, strings.ReplaceAll(skillIDLower, "_", " ")) {
		score += 4.5
		exact = true
		matchedTags = append(matchedTags, "exact_id")
	}

	// Check trigger matches with word boundary
	for _, tr := range skill.Triggers {
		if matchWordBoundary(rawPromptLower, tr) {
			score += 3.5
			matchedTags = append(matchedTags, "trigger:"+tr)
		}
	}

	// Check tool name matches with word boundary
	for _, tool := range skill.Tools {
		toolLower := strings.ToLower(tool.Name)
		if matchWordBoundary(rawPromptLower, toolLower) || matchWordBoundary(rawPromptLower, strings.ReplaceAll(toolLower, "_", " ")) {
			score += 4.0
			matchedTags = append(matchedTags, "tool:"+tool.Name)
		}
	}

	// Check keyword overlap
	keywordMap := make(map[string]bool)
	for _, kw := range skill.Keywords {
		keywordMap[strings.ToLower(kw)] = true
	}

	for _, token := range queryTokens {
		if keywordMap[token] {
			score += 0.3
		}
	}

	// Check domain and category match
	if skill.Domain != "" && matchWordBoundary(rawPromptLower, skill.Domain) {
		score += 1.0
	}
	if skill.Category != "" && matchWordBoundary(rawPromptLower, skill.Category) {
		score += 0.8
	}

	// Require at least one strong anchor (name, id, trigger, tool) to exceed refusal threshold
	if !exact && len(matchedTags) == 0 {
		if score > 1.0 {
			score = 1.0
		}
	}

	return score, exact, matchedTags
}

func matchWordBoundary(text, target string) bool {
	t := strings.TrimSpace(strings.ToLower(target))
	if t == "" {
		return false
	}
	pattern := `(?i)(?:^|[^a-zA-Z0-9_])` + regexp.QuoteMeta(t) + `(?:$|[^a-zA-Z0-9_])`
	matched, err := regexp.MatchString(pattern, text)
	return err == nil && matched
}

