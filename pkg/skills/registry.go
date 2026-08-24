package skills

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// SkillRegistry holds all active skills in memory with fast thread-safe indexing
type SkillRegistry struct {
	mu           sync.RWMutex
	skills       map[string]*Skill
	skillsByName map[string]*Skill
	toolsByName  map[string]ToolDefinition
	toolToSkill  map[string]*Skill
	byCategory   map[string][]*Skill
	parser       *Parser
	globalSeal   string
}

// NewSkillRegistry creates an initialized SkillRegistry
func NewSkillRegistry() *SkillRegistry {
	return &SkillRegistry{
		skills:       make(map[string]*Skill),
		skillsByName: make(map[string]*Skill),
		toolsByName:  make(map[string]ToolDefinition),
		toolToSkill:  make(map[string]*Skill),
		byCategory:   make(map[string][]*Skill),
		parser:       NewParser(),
	}
}

// Register adds or updates a skill in the registry
func (r *SkillRegistry) Register(skill *Skill) {
	if skill == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.skills[skill.ID] = skill
	r.skillsByName[strings.ToLower(skill.Name)] = skill
	for _, tool := range skill.Tools {
		r.toolsByName[tool.Name] = tool
		r.toolToSkill[tool.Name] = skill
	}

	if skill.Category != "" {
		catLower := strings.ToLower(skill.Category)
		r.byCategory[catLower] = append(r.byCategory[catLower], skill)
	}

	r.recomputeGlobalSeal()
}

// Unregister removes a skill from the registry by ID
func (r *SkillRegistry) Unregister(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	skill, exists := r.skills[id]
	if !exists {
		return false
	}

	delete(r.skills, id)
	delete(r.skillsByName, strings.ToLower(skill.Name))
	for _, tool := range skill.Tools {
		delete(r.toolsByName, tool.Name)
		delete(r.toolToSkill, tool.Name)
	}

	if skill.Category != "" {
		catLower := strings.ToLower(skill.Category)
		list := r.byCategory[catLower]
		filtered := make([]*Skill, 0, len(list))
		for _, s := range list {
			if s.ID != id {
				filtered = append(filtered, s)
			}
		}
		r.byCategory[catLower] = filtered
	}

	r.recomputeGlobalSeal()
	return true
}

// LoadFromDirectory recursively scans and loads all .md and .json skills from a folder
func (r *SkillRegistry) LoadFromDirectory(dirPath string) (int, error) {
	count := 0
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" || ext == ".json" {
			skill, parseErr := r.parser.ParseFile(path)
			if parseErr == nil && skill != nil {
				r.Register(skill)
				count++
			}
		}
		return nil
	})

	return count, err
}

// GetByID returns a skill by its exact ID
func (r *SkillRegistry) GetByID(id string) (*Skill, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	skill, exists := r.skills[id]
	return skill, exists
}

// GetByName returns a skill by its human-readable name (case-insensitive)
func (r *SkillRegistry) GetByName(name string) (*Skill, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	skill, exists := r.skillsByName[strings.ToLower(name)]
	return skill, exists
}

// GetToolByName returns a tool and its parent skill by tool name
func (r *SkillRegistry) GetToolByName(toolName string) (ToolDefinition, *Skill, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, exists := r.toolsByName[toolName]
	skill := r.toolToSkill[toolName]
	return tool, skill, exists
}

// GetByCategory returns all skills belonging to a category
func (r *SkillRegistry) GetByCategory(category string) []*Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.byCategory[strings.ToLower(category)]
	result := make([]*Skill, len(list))
	copy(result, list)
	return result
}

// GetAllSkills returns a slice of all registered skills sorted by ID
func (r *SkillRegistry) GetAllSkills() []*Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]*Skill, 0, len(r.skills))
	for _, s := range r.skills {
		list = append(list, s)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})
	return list
}

// GetAllTools returns all tools across all skills sorted by name
func (r *SkillRegistry) GetAllTools() []ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tools := make([]ToolDefinition, 0, len(r.toolsByName))
	for _, t := range r.toolsByName {
		tools = append(tools, t)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})
	return tools
}

// Count returns the total number of registered skills
func (r *SkillRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.skills)
}

// TotalMonolithicTokens computes the total prompt tokens if ALL skills were injected in-context
func (r *SkillRegistry) TotalMonolithicTokens() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, s := range r.skills {
		total += s.EstimatedTokenSize()
	}
	return total
}

// ExportMCPTools converts all registered tools into standard MCP tools/list format
func (r *SkillRegistry) ExportMCPTools() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mcpList := MCPListToolsResponse{
		Tools: make([]MCPToolDefinition, 0, len(r.toolsByName)),
	}
	for _, t := range r.toolsByName {
		mcpList.Tools = append(mcpList.Tools, MCPToolDefinition{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.Parameters,
		})
	}
	return json.MarshalIndent(mcpList, "", "  ")
}

// ExportJSON exports the full registry to a JSON payload
func (r *SkillRegistry) ExportJSON() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return json.MarshalIndent(r.GetAllSkills(), "", "  ")
}

// GlobalMerkleSeal returns the combined Merkle state hash of the entire registry
func (r *SkillRegistry) GlobalMerkleSeal() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.globalSeal
}

func (r *SkillRegistry) recomputeGlobalSeal() {
	h := sha256.New()
	sortedIDs := make([]string, 0, len(r.skills))
	for id := range r.skills {
		sortedIDs = append(sortedIDs, id)
	}
	sort.Strings(sortedIDs)

	for _, id := range sortedIDs {
		h.Write([]byte(r.skills[id].MerkleSeal))
	}
	r.globalSeal = hex.EncodeToString(h.Sum(nil))
}

// Summary returns a diagnostic string of loaded skills
func (r *SkillRegistry) Summary() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sealPrefix := ""
	if len(r.globalSeal) >= 12 {
		sealPrefix = r.globalSeal[:12]
	}
	return fmt.Sprintf("SkillRegistry: %d skills loaded | %d tools available | %d estimated monolithic tokens | Merkle: %s",
		len(r.skills), len(r.toolsByName), r.TotalMonolithicTokens(), sealPrefix)
}
