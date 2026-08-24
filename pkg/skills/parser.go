package skills

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Parser handles parsing markdown skills and MCP configurations
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// ParseMarkdownSkill parses a standard SKILL.md or markdown skill specification
func (p *Parser) ParseMarkdownSkill(content string, sourcePath string) (*Skill, error) {
	skill := &Skill{
		CreatedAt:  time.Now(),
		SourcePath: sourcePath,
	}

	// 1. Extract Frontmatter if present
	frontmatterRegex := regexp.MustCompile(`(?s)^---\r?\n(.*?)\r?\n---`)
	matches := frontmatterRegex.FindStringSubmatch(content)
	body := content

	if len(matches) > 1 {
		fmContent := matches[1]
		body = strings.TrimPrefix(content, matches[0])
		p.parseYAMLFrontmatter(fmContent, skill)
	}

	// If name not found in frontmatter, parse primary H1 header
	if skill.Name == "" {
		h1Regex := regexp.MustCompile(`(?m)^#\s+(.+)$`)
		if h1Match := h1Regex.FindStringSubmatch(body); len(h1Match) > 1 {
			skill.Name = strings.TrimSpace(h1Match[1])
		} else {
			base := filepath.Base(sourcePath)
			skill.Name = strings.TrimSuffix(base, filepath.Ext(base))
		}
	}

	if skill.ID == "" {
		skill.ID = strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(skill.Name, "_"))
	}

	// 2. Extract Description
	if skill.Description == "" {
		descRegex := regexp.MustCompile(`(?m)^>\s*(.+)$`)
		if descMatch := descRegex.FindStringSubmatch(body); len(descMatch) > 1 {
			skill.Description = strings.TrimSpace(descMatch[1])
		} else {
			lines := strings.Split(body, "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
					skill.Description = trimmed
					break
				}
			}
		}
	}

	// 3. Extract Tools from markdown codeblocks or section
	p.extractToolsFromMarkdown(body, skill)

	// 4. Instructions
	skill.Instructions = strings.TrimSpace(body)

	// 5. Extract Keywords and Triggers if not present
	if len(skill.Triggers) == 0 {
		skill.Triggers = p.extractTriggers(skill.Name, skill.Description)
	}
	if len(skill.Keywords) == 0 {
		skill.Keywords = p.extractKeywords(skill.Name, skill.Description, skill.Instructions)
	}

	// 6. Compute SHA-256 content hash (field name MerkleSeal is historical)
	h := sha256.New()
	h.Write([]byte(skill.ID + ":" + skill.Name + ":" + skill.Instructions))
	for _, t := range skill.Tools {
		toolBytes, _ := json.Marshal(t)
		h.Write(toolBytes)
	}
	skill.MerkleSeal = hex.EncodeToString(h.Sum(nil))

	return skill, nil
}

func (p *Parser) parseYAMLFrontmatter(fm string, skill *Skill) {
	lines := strings.Split(fm, "\n")
	currentListKey := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for YAML list items: "- item"
		if strings.HasPrefix(line, "- ") && currentListKey != "" {
			item := strings.Trim(strings.TrimPrefix(line, "- "), "\"' \t")
			if item != "" {
				switch currentListKey {
				case "triggers":
					skill.Triggers = append(skill.Triggers, item)
				case "keywords":
					skill.Keywords = append(skill.Keywords, item)
				case "examples":
					skill.Examples = append(skill.Examples, item)
				}
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

		currentListKey = key

		switch key {
		case "name":
			skill.Name = val
		case "id":
			skill.ID = val
		case "version":
			skill.Version = val
		case "description":
			skill.Description = val
		case "category":
			skill.Category = val
		case "domain":
			skill.Domain = val
		case "author":
			skill.Author = val
		case "triggers":
			if val != "" {
				var trList []string
				if err := json.Unmarshal([]byte(strings.TrimSpace(parts[1])), &trList); err == nil {
					skill.Triggers = trList
				} else if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
					rawList := strings.Trim(val, "[]")
					for _, item := range strings.Split(rawList, ",") {
						t := strings.Trim(strings.TrimSpace(item), "\"'")
						if t != "" {
							skill.Triggers = append(skill.Triggers, t)
						}
					}
				} else {
					skill.Triggers = append(skill.Triggers, val)
				}
			}
		case "keywords":
			if val != "" {
				var kwList []string
				if err := json.Unmarshal([]byte(strings.TrimSpace(parts[1])), &kwList); err == nil {
					skill.Keywords = kwList
				} else if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
					rawList := strings.Trim(val, "[]")
					for _, item := range strings.Split(rawList, ",") {
						t := strings.Trim(strings.TrimSpace(item), "\"'")
						if t != "" {
							skill.Keywords = append(skill.Keywords, t)
						}
					}
				} else {
					skill.Keywords = append(skill.Keywords, val)
				}
			}
		}
	}
}

func (p *Parser) extractToolsFromMarkdown(body string, skill *Skill) {
	// Look for json code blocks that define tools
	jsonBlockRegex := regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\}|\\[.*?\\])\\s*```")
	matches := jsonBlockRegex.FindAllStringSubmatch(body, -1)

	for _, m := range matches {
		if len(m) > 1 {
			rawJSON := strings.TrimSpace(m[1])

			// 1. Try MCP format: {"inputSchema": ...}
			var mcpTool MCPToolDefinition
			if err := json.Unmarshal([]byte(rawJSON), &mcpTool); err == nil && mcpTool.Name != "" && mcpTool.InputSchema.Type != "" {
				skill.Tools = append(skill.Tools, ToolDefinition{
					Name:        mcpTool.Name,
					Description: mcpTool.Description,
					Parameters:  mcpTool.InputSchema,
					Category:    skill.Category,
				})
				continue
			}

			// 2. Try standard ToolDefinition
			var singleTool ToolDefinition
			if err := json.Unmarshal([]byte(rawJSON), &singleTool); err == nil && singleTool.Name != "" {
				if singleTool.Category == "" {
					singleTool.Category = skill.Category
				}
				skill.Tools = append(skill.Tools, singleTool)
				continue
			}

			// 3. Try slice of ToolDefinitions
			var multiTools []ToolDefinition
			if err := json.Unmarshal([]byte(rawJSON), &multiTools); err == nil && len(multiTools) > 0 {
				for i := range multiTools {
					if multiTools[i].Category == "" {
						multiTools[i].Category = skill.Category
					}
				}
				skill.Tools = append(skill.Tools, multiTools...)
				continue
			}

			// 4. Try MCP list tools response wrapper {"tools": [...]}
			var mcpWrapper struct {
				Tools []MCPToolDefinition `json:"tools"`
			}
			if err := json.Unmarshal([]byte(rawJSON), &mcpWrapper); err == nil && len(mcpWrapper.Tools) > 0 {
				for _, mt := range mcpWrapper.Tools {
					skill.Tools = append(skill.Tools, ToolDefinition{
						Name:        mt.Name,
						Description: mt.Description,
						Parameters:  mt.InputSchema,
						Category:    skill.Category,
					})
				}
			}
		}
	}

	// If no tools explicitly defined via JSON, synthesize a default canonical tool for this skill
	if len(skill.Tools) == 0 {
		toolName := strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(skill.Name, "_"))
		skill.Tools = append(skill.Tools, ToolDefinition{
			Name:        toolName,
			Description: fmt.Sprintf("Execute action for skill: %s. %s", skill.Name, skill.Description),
			Parameters: ToolParameters{
				Type: "object",
				Properties: map[string]ParameterProperty{
					"query": {
						Type:        "string",
						Description: "The query or instruction to execute for this tool",
					},
					"params": {
						Type:        "string",
						Description: "Optional JSON parameters for the execution",
					},
				},
				Required: []string{"query"},
			},
			Category: skill.Category,
		})
	}
}

func (p *Parser) extractTriggers(name, desc string) []string {
	words := strings.Fields(strings.ToLower(name + " " + desc))
	triggers := make([]string, 0)
	seen := make(map[string]bool)

	for _, w := range words {
		cleaned := regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(w, "")
		if len(cleaned) > 3 && !seen[cleaned] {
			seen[cleaned] = true
			triggers = append(triggers, cleaned)
		}
	}
	return triggers
}

func (p *Parser) extractKeywords(name, desc, instructions string) []string {
	combined := strings.ToLower(name + " " + desc + " " + instructions)
	words := strings.Fields(combined)
	keywords := make([]string, 0)
	seen := make(map[string]bool)

	stopwords := map[string]bool{
		"this": true, "that": true, "with": true, "from": true, "have": true,
		"will": true, "your": true, "when": true, "what": true, "then": true,
		"into": true, "more": true, "some": true, "must": true, "only": true,
		"about": true, "after": true, "also": true, "could": true, "which": true,
	}

	for _, w := range words {
		cleaned := regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(w, "")
		if len(cleaned) > 3 && !seen[cleaned] && !stopwords[cleaned] {
			seen[cleaned] = true
			keywords = append(keywords, cleaned)
			if len(keywords) >= 25 {
				break
			}
		}
	}
	return keywords
}

// ParseFile parses a file from disk
func (p *Parser) ParseFile(filePath string) (*Skill, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill file %s: %w", filePath, err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".json" {
		var skill Skill
		if err := json.Unmarshal(bytes, &skill); err == nil && skill.Name != "" {
			skill.SourcePath = filePath
			if skill.ID == "" {
				skill.ID = strings.ToLower(regexp.MustCompile(`[^a-zA-Z0-9_]+`).ReplaceAllString(skill.Name, "_"))
			}
			return &skill, nil
		}
	}

	return p.ParseMarkdownSkill(string(bytes), filePath)
}
