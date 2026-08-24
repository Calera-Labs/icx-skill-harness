package skills

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ParameterProperty describes a single parameter property in a tool schema (JSON Schema format)
type ParameterProperty struct {
	Type        string                       `json:"type"`
	Description string                       `json:"description,omitempty"`
	Enum        []string                     `json:"enum,omitempty"`
	Default     any                          `json:"default,omitempty"`
	Items       *ParameterProperty           `json:"items,omitempty"`
	Properties  map[string]ParameterProperty `json:"properties,omitempty"`
}

// ToolParameters describes the JSON schema parameters for a tool
type ToolParameters struct {
	Type       string                       `json:"type"`
	Properties map[string]ParameterProperty `json:"properties"`
	Required   []string                     `json:"required,omitempty"`
}

// ToolDefinition represents a single callable function/tool associated with a skill
type ToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
	Category    string         `json:"category,omitempty"`
	Strict      bool           `json:"strict,omitempty"`
}

// Validate checks if the tool definition has valid names and schema
func (t *ToolDefinition) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("tool description cannot be empty for %s", t.Name)
	}
	if t.Parameters.Type == "" {
		t.Parameters.Type = "object"
	}
	if t.Parameters.Properties == nil {
		t.Parameters.Properties = make(map[string]ParameterProperty)
	}
	return nil
}

// Skill represents a full skill definition (parsed from SKILL.md or MCP)
type Skill struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Version      string           `json:"version,omitempty"`
	Description  string           `json:"description"`
	Category     string           `json:"category,omitempty"`
	Domain       string           `json:"domain,omitempty"`
	Author       string           `json:"author,omitempty"`
	Triggers     []string         `json:"triggers"`
	Keywords     []string         `json:"keywords"`
	Tools        []ToolDefinition `json:"tools"`
	Instructions string           `json:"instructions"`
	Examples     []string         `json:"examples,omitempty"`
	SourcePath   string           `json:"source_path,omitempty"`
	MerkleSeal   string           `json:"merkle_seal"` // SHA-256 content hash; not a Merkle tree
	CreatedAt    time.Time        `json:"created_at"`
}

// Validate checks if the skill definition is well-formed
func (s *Skill) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("skill ID cannot be empty")
	}
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("skill name cannot be empty")
	}
	for i := range s.Tools {
		if err := s.Tools[i].Validate(); err != nil {
			return fmt.Errorf("skill %s has invalid tool: %w", s.ID, err)
		}
	}
	return nil
}

// EstimatedTokenSize calculates an approximate token count for this skill's schema and instructions
func (s *Skill) EstimatedTokenSize() int {
	data, _ := json.Marshal(s.Tools)
	// Rough estimate: ~4 chars per token for JSON schema and instructions + metadata overhead
	schemaTokens := len(data) / 4
	instrTokens := len(s.Instructions) / 4
	return schemaTokens + instrTokens + 50
}

// SkillViewport represents the extracted micro-viewport injected into the LLM
type SkillViewport struct {
	MatchedSkills     []*Skill         `json:"matched_skills"`
	ActiveTools       []ToolDefinition `json:"active_tools"`
	TotalSchemaTokens int              `json:"total_schema_tokens"`
	ConfidenceScore   float64          `json:"confidence_score"`
	QueryIntent       string           `json:"query_intent"`
	IsRefusal         bool             `json:"is_refusal"`
	RefusalReason     string           `json:"refusal_reason,omitempty"`
	RoutingLatencyUs  int64            `json:"routing_latency_us"`
	LatticeWalkUsed   bool             `json:"lattice_walk_used,omitempty"`
}

// MCPToolDefinition represents standard MCP format tool description
type MCPToolDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema ToolParameters `json:"inputSchema"`
}

// MCPListToolsResponse represents a standard MCP tools/list response
type MCPListToolsResponse struct {
	Tools []MCPToolDefinition `json:"tools"`
}
