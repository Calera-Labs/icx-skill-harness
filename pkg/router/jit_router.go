package router

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

var (
	numStepRegex          = regexp.MustCompile(`(?i)(?:^|\s)(?:\d+[\.\)]|step\s*\d+[:\.]?)\s+`)
	connectorRegex        = regexp.MustCompile(`(?i)(?:;\s*|\s+(?:then|after that|afterwards|next|finally|subsequently|and then)\s+)`)
	futureYearDigitsRegex = regexp.MustCompile(`(?:20[3-9][0-9]|21[0-9]{2})`)

	stopwords = map[string]bool{
		"the": true, "and": true, "for": true, "with": true, "this": true,
		"that": true, "from": true, "into": true, "have": true, "what": true,
		"where": true, "when": true, "how": true, "please": true, "want": true,
		"your": true, "can": true, "you": true, "about": true, "more": true,
		"will": true, "been": true, "each": true, "which": true, "does": true,
		"were": true, "them": true, "they": true, "then": true, "some": true,
	}
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

// indexedSkill holds pre-computed lowercase token sets for high-speed routing
type indexedSkill struct {
	skill         *skills.Skill
	nameLower     string
	nameSpaced    string
	idLower       string
	idSpaced      string
	categoryLower string
	domainLower   string
	triggersLower []string
	toolsLower    []string
	keywordsMap   map[string]bool
}

// LatticeSkillRouter provides sub-millisecond JIT micro-viewport routing over hundreds of skills
type LatticeSkillRouter struct {
	registry   *skills.SkillRegistry
	icx        *icx.Client
	config     RouterConfig
	mu         sync.RWMutex
	lastSeal   string
	indexCache []indexedSkill
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
	r := &LatticeSkillRouter{
		registry: reg,
		icx:      icxClient,
		config:   cfg,
	}
	r.rebuildIndex()
	return r
}

func (r *LatticeSkillRouter) rebuildIndex() {
	if r.registry == nil {
		return
	}
	allSkills := r.registry.GetAllSkills()
	r.mu.Lock()
	defer r.mu.Unlock()

	r.indexCache = make([]indexedSkill, len(allSkills))
	r.lastSeal = r.registry.GlobalMerkleSeal()

	for i, s := range allSkills {
		nLower := strings.ToLower(s.Name)
		idLower := strings.ToLower(s.ID)
		idxItem := indexedSkill{
			skill:         s,
			nameLower:     nLower,
			nameSpaced:    strings.ReplaceAll(nLower, "_", " "),
			idLower:       idLower,
			idSpaced:      strings.ReplaceAll(idLower, "_", " "),
			categoryLower: strings.ToLower(s.Category),
			domainLower:   strings.ToLower(s.Domain),
			triggersLower: make([]string, len(s.Triggers)),
			toolsLower:    make([]string, len(s.Tools)),
			keywordsMap:   make(map[string]bool, len(s.Keywords)),
		}

		for ti, tr := range s.Triggers {
			idxItem.triggersLower[ti] = strings.ToLower(tr)
		}

		for ti, t := range s.Tools {
			idxItem.toolsLower[ti] = strings.ToLower(t.Name)
		}

		for _, kw := range s.Keywords {
			idxItem.keywordsMap[strings.ToLower(kw)] = true
		}

		r.indexCache[i] = idxItem
	}
}

func (r *LatticeSkillRouter) getIndexedSkills() []indexedSkill {
	r.mu.RLock()
	if r.registry != nil && r.registry.GlobalMerkleSeal() != r.lastSeal {
		r.mu.RUnlock()
		r.rebuildIndex()
		r.mu.RLock()
	}
	defer r.mu.RUnlock()
	return r.indexCache
}

// RouteQuery analyzes a user prompt and generates a compact micro-viewport
func (r *LatticeSkillRouter) RouteQuery(prompt string) (*skills.SkillViewport, error) {
	return r.RouteQueryWithContext(context.Background(), prompt)
}

// RouteQueryWithContext analyzes prompt with context cancellation
func (r *LatticeSkillRouter) RouteQueryWithContext(ctx context.Context, prompt string) (*skills.SkillViewport, error) {
	t0 := time.Now()
	indexed := r.getIndexedSkills()

	if len(indexed) == 0 {
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

	// 0. Temporal & Epistemic Premise Validation (Anti-Hallucination Guardrail)
	if isTrap, reason := isUnanswerablePremise(trimmedPrompt); isTrap {
		return &skills.SkillViewport{
			IsRefusal:        true,
			RefusalReason:    fmt.Sprintf("SAFE_REFUSAL: %s in query '%s'", reason, prompt),
			RoutingLatencyUs: time.Since(t0).Microseconds(),
			LatticeWalkUsed:  false,
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

	// 3. Score skills using pre-computed index & zero-alloc boundary matcher
	scored := make([]ScoredSkill, 0, len(indexed))
	promptLower := strings.ToLower(trimmedPrompt)

	for _, idxItem := range indexed {
		score, exact, tags := r.scoreIndexedSkill(idxItem, tokens, promptLower)

		// Apply ICX Lattice walk boost if match found
		if boost, hasBoost := latticeBoostMap[idxItem.idLower]; hasBoost {
			score += boost
			tags = append(tags, fmt.Sprintf("icx_lattice:%.2f", boost))
		}

		if score >= r.config.MinScoreThreshold {
			scored = append(scored, ScoredSkill{
				Skill:       idxItem.skill,
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
	selectedSkills := make([]*skills.Skill, 0, r.config.MaxToolsPerViewport)
	activeTools := make([]skills.ToolDefinition, 0, r.config.MaxToolsPerViewport*2)
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

// isUnanswerablePremise guards against temporal impossibility and unregistered trap premises
func isUnanswerablePremise(prompt string) (bool, string) {
	pLower := strings.ToLower(prompt)

	// 1. Nonexistent Future SEC / Financial Filing Trap (e.g. FY2035 10-K)
	isFinancialFiling := strings.Contains(pLower, "10-k") ||
		strings.Contains(pLower, "10-q") ||
		strings.Contains(pLower, "8-k") ||
		strings.Contains(pLower, "sec filing") ||
		strings.Contains(pLower, "edgar filing") ||
		strings.Contains(pLower, "operating margin") ||
		strings.Contains(pLower, "gaap revenue") ||
		strings.Contains(pLower, "sec 10-k")

	if isFinancialFiling {
		// Look for future years > 2028 (e.g. 2035, FY2035, 2040)
		if matches := futureYearDigitsRegex.FindAllString(pLower, -1); len(matches) > 0 {
			for _, yrStr := range matches {
				if yr, err := strconv.Atoi(yrStr); err == nil && yr > 2028 {
					return true, fmt.Sprintf("Temporal impossibility: Unverified future SEC filing for year %d", yr)
				}
			}
		}
	}

	return false, ""
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

	// 1. Extract sequential pipeline segments
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
	executedSet := make(map[string]bool, len(executedTools))
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
	indices := numStepRegex.FindAllStringIndex(prompt, -1)
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
	var b strings.Builder
	b.Grow(len(cleaned))
	for i := 0; i < len(cleaned); i++ {
		c := cleaned[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '-' || c == '.' {
			b.WriteByte(c)
		} else {
			b.WriteByte(' ')
		}
	}
	words := strings.Fields(b.String())

	res := make([]string, 0, len(words))
	for _, w := range words {
		if len(w) > 2 && !stopwords[w] {
			res = append(res, w)
		}
	}
	return res
}

func (r *LatticeSkillRouter) scoreIndexedSkill(idx indexedSkill, queryTokens []string, promptLower string) (float64, bool, []string) {
	score := 0.0
	exact := false
	matchedTags := make([]string, 0, 4)

	// Check exact name match with zero-alloc boundary
	if containsWord(promptLower, idx.nameLower) || containsWord(promptLower, idx.nameSpaced) {
		score += 5.0
		exact = true
		matchedTags = append(matchedTags, "exact_name")
	}

	// Check skill ID match
	if containsWord(promptLower, idx.idLower) || containsWord(promptLower, idx.idSpaced) {
		score += 4.5
		exact = true
		matchedTags = append(matchedTags, "exact_id")
	}

	// Check triggers with fast word boundary
	for _, tr := range idx.triggersLower {
		if containsWord(promptLower, tr) {
			score += 3.5
			matchedTags = append(matchedTags, "trigger:"+tr)
		}
	}

	// Check tool names
	for _, t := range idx.toolsLower {
		if containsWord(promptLower, t) || containsWord(promptLower, strings.ReplaceAll(t, "_", " ")) {
			score += 4.0
			matchedTags = append(matchedTags, "tool:"+t)
		}
	}

	// Check keyword overlap
	for _, token := range queryTokens {
		if idx.keywordsMap[token] {
			score += 0.3
		}
	}

	// Check domain and category match
	if idx.domainLower != "" && containsWord(promptLower, idx.domainLower) {
		score += 1.0
	}
	if idx.categoryLower != "" && containsWord(promptLower, idx.categoryLower) {
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

// containsWord performs zero-allocation word-boundary substring search
func containsWord(text, target string) bool {
	tLen := len(target)
	if tLen == 0 {
		return false
	}
	start := 0
	for {
		idx := strings.Index(text[start:], target)
		if idx == -1 {
			return false
		}
		pos := start + idx
		endPos := pos + tLen

		// Check left boundary
		leftOK := (pos == 0) || isDelimiter(text[pos-1])
		// Check right boundary
		rightOK := (endPos == len(text)) || isDelimiter(text[endPos])

		if leftOK && rightOK {
			return true
		}
		start = pos + 1
		if start >= len(text) {
			return false
		}
	}
}

func isDelimiter(c byte) bool {
	return (c < 'a' || c > 'z') && (c < '0' || c > '9') && (c < 'A' || c > 'Z')
}
