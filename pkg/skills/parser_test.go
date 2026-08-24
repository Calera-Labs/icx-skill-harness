package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMarkdownSkill(t *testing.T) {
	parser := NewParser()

	content := `---
name: "PostgreSQL Database Administrator"
id: "postgres_db_admin"
category: "Database"
description: "Executes parameterized SQL queries and schema migrations"
triggers: ["postgres", "sql", "database", "transactions table"]
keywords: ["postgres", "query", "table", "migration", "schema"]
---

# PostgreSQL Database Administrator

Executes parameterized SQL queries against PostgreSQL.

` + "```json" + `
{
  "name": "postgres_executor",
  "description": "Execute SQL queries against PostgreSQL database",
  "parameters": {
    "type": "object",
    "properties": {
      "query": {"type": "string", "description": "SQL query to execute"},
      "options": {"type": "string", "description": "Optional parameters"}
    },
    "required": ["query"]
  }
}
` + "```"

	skill, err := parser.ParseMarkdownSkill(content, "test.md")
	if err != nil {
		t.Fatalf("unexpected error parsing markdown: %v", err)
	}

	if skill.ID != "postgres_db_admin" {
		t.Errorf("expected ID 'postgres_db_admin', got '%s'", skill.ID)
	}
	if skill.Name != "PostgreSQL Database Administrator" {
		t.Errorf("expected Name 'PostgreSQL Database Administrator', got '%s'", skill.Name)
	}
	if skill.Category != "Database" {
		t.Errorf("expected Category 'Database', got '%s'", skill.Category)
	}
	if len(skill.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(skill.Tools))
	}
	if skill.Tools[0].Name != "postgres_executor" {
		t.Errorf("expected tool name 'postgres_executor', got '%s'", skill.Tools[0].Name)
	}
	if skill.MerkleSeal == "" {
		t.Errorf("expected non-empty MerkleSeal")
	}
}

func TestParseMarkdownWithoutFrontmatter(t *testing.T) {
	parser := NewParser()

	content := `# Custom Tool Runner
> A runner for custom shell operations

Executes custom shell operations safely.
`
	skill, err := parser.ParseMarkdownSkill(content, "custom_runner.md")
	if err != nil {
		t.Fatalf("unexpected error parsing markdown: %v", err)
	}

	if skill.Name != "Custom Tool Runner" {
		t.Errorf("expected name 'Custom Tool Runner', got '%s'", skill.Name)
	}
	if len(skill.Tools) != 1 {
		t.Fatalf("expected 1 synthesized tool, got %d", len(skill.Tools))
	}
	if skill.Tools[0].Name != "custom_tool_runner" {
		t.Errorf("expected synthesized tool name 'custom_tool_runner', got '%s'", skill.Tools[0].Name)
	}
}

func TestParseMCPToolFormat(t *testing.T) {
	parser := NewParser()

	content := `---
name: "Weather Reporter"
---

# Weather Reporter

` + "```json" + `
{
  "name": "get_weather",
  "description": "Get current weather for location",
  "inputSchema": {
    "type": "object",
    "properties": {
      "location": {"type": "string", "description": "City name"}
    },
    "required": ["location"]
  }
}
` + "```"

	skill, err := parser.ParseMarkdownSkill(content, "weather.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(skill.Tools) != 1 {
		t.Fatalf("expected 1 MCP tool, got %d", len(skill.Tools))
	}
	if skill.Tools[0].Name != "get_weather" {
		t.Errorf("expected tool name 'get_weather', got '%s'", skill.Tools[0].Name)
	}
	if _, ok := skill.Tools[0].Parameters.Properties["location"]; !ok {
		t.Errorf("expected 'location' property in input schema")
	}
}

func TestParseJSONFile(t *testing.T) {
	parser := NewParser()
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "sample_skill.json")

	jsonContent := `{
		"id": "json_skill_01",
		"name": "JSON Configured Skill",
		"description": "Skill loaded directly from JSON",
		"category": "Testing",
		"triggers": ["json", "config"],
		"tools": [
			{
				"name": "json_exec",
				"description": "Execute JSON task",
				"parameters": {
					"type": "object",
					"properties": {
						"task": {"type": "string"}
					},
					"required": ["task"]
				}
			}
		]
	}`

	if err := os.WriteFile(jsonPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	skill, err := parser.ParseFile(jsonPath)
	if err != nil {
		t.Fatalf("unexpected error parsing JSON file: %v", err)
	}

	if skill.ID != "json_skill_01" {
		t.Errorf("expected ID 'json_skill_01', got '%s'", skill.ID)
	}
	if skill.Name != "JSON Configured Skill" {
		t.Errorf("expected Name 'JSON Configured Skill', got '%s'", skill.Name)
	}
}
