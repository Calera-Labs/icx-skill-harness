package skills

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// OpenAIToolParamSchema represents the parameters object in OpenAI tool schema
type OpenAIToolParamSchema struct {
	Type       string                       `json:"type"`
	Properties map[string]ParameterProperty `json:"properties"`
	Required   []string                     `json:"required,omitempty"`
}

// ConvertOpenAIToolToDefinition converts an OpenAI-format tool into an internal ToolDefinition
func ConvertOpenAIToolToDefinition(name, description string, params map[string]any) ToolDefinition {
	paramBytes, _ := json.Marshal(params)
	var toolParams ToolParameters
	_ = json.Unmarshal(paramBytes, &toolParams)

	if toolParams.Type == "" {
		toolParams.Type = "object"
	}
	if toolParams.Properties == nil {
		toolParams.Properties = make(map[string]ParameterProperty)
	}

	return ToolDefinition{
		Name:        name,
		Description: description,
		Parameters:  toolParams,
		Category:    inferCategoryFromText(name + " " + description),
	}
}

// CreateSkillFromTool builds a standalone Skill wrapper around a dynamic tool definition
func CreateSkillFromTool(tool ToolDefinition) *Skill {
	id := strings.ToLower(tool.Name)
	triggers := ExtractTriggers(tool.Name, tool.Description)

	h := sha256.New()
	h.Write([]byte(tool.Name + ":" + tool.Description))
	seal := hex.EncodeToString(h.Sum(nil))

	return &Skill{
		ID:           id,
		Name:         tool.Name,
		Description:  tool.Description,
		Category:     tool.Category,
		Triggers:     triggers,
		Keywords:     triggers,
		Tools:        []ToolDefinition{tool},
		Instructions: fmt.Sprintf("Dynamic tool '%s': %s", tool.Name, tool.Description),
		MerkleSeal:   seal,
		CreatedAt:    time.Now(),
	}
}

// ExtractTriggers tokenizes tool name and description into keyword trigger phrases
func ExtractTriggers(name, description string) []string {
	combined := strings.ToLower(name + " " + description)
	combined = strings.ReplaceAll(combined, "_", " ")
	combined = strings.ReplaceAll(combined, "-", " ")
	combined = strings.ReplaceAll(combined, ".", " ")

	re := regexp.MustCompile(`[a-z0-9]{3,}`)
	words := re.FindAllString(combined, -1)

	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "with": true, "from": true,
		"this": true, "that": true, "tool": true, "query": true, "exec": true,
		"primary": true, "action": true, "optional": true, "parameters": true,
		"string": true, "object": true, "description": true,
	}

	seen := make(map[string]bool)
	var triggers []string

	// Always add the base tool name parts
	nameClean := strings.ToLower(strings.ReplaceAll(name, "_", " "))
	if nameClean != "" && !seen[nameClean] {
		seen[nameClean] = true
		triggers = append(triggers, nameClean)
	}

	for _, w := range words {
		if !stopWords[w] && !seen[w] {
			seen[w] = true
			triggers = append(triggers, w)
		}
	}

	return triggers
}

func inferCategoryFromText(text string) string {
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "sql") || strings.Contains(lower, "db") || strings.Contains(lower, "postgres") || strings.Contains(lower, "database"):
		return "Database"
	case strings.Contains(lower, "sec") || strings.Contains(lower, "finance") || strings.Contains(lower, "valuation") || strings.Contains(lower, "10-k") || strings.Contains(lower, "stock"):
		return "Finance"
	case strings.Contains(lower, "k8s") || strings.Contains(lower, "docker") || strings.Contains(lower, "cloud") || strings.Contains(lower, "deploy") || strings.Contains(lower, "infra"):
		return "DevOps"
	case strings.Contains(lower, "protein") || strings.Contains(lower, "gene") || strings.Contains(lower, "variant") || strings.Contains(lower, "clinical"):
		return "Life Sciences"
	case strings.Contains(lower, "git") || strings.Contains(lower, "code") || strings.Contains(lower, "ast") || strings.Contains(lower, "repo"):
		return "Developer Tools"
	default:
		return "Custom Tools"
	}
}
