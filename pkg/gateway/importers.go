package gateway

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

// MCPConfigFile represents the standard Claude Desktop / Cursor MCP JSON format
type MCPConfigFile struct {
	MCPServers map[string]MCPServerEntry `json:"mcpServers"`
}

// MCPServerEntry represents a single MCP server configuration
type MCPServerEntry struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// LoadMCPServersConfig parses an MCP JSON config file and registers synthetic or configured tool definitions
func LoadMCPServersConfig(filePath string, reg *skills.SkillRegistry) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read MCP config: %w", err)
	}

	var mcpConfig MCPConfigFile
	if err := json.Unmarshal(data, &mcpConfig); err != nil {
		return 0, fmt.Errorf("failed to parse MCP config JSON: %w", err)
	}

	count := 0
	for serverName := range mcpConfig.MCPServers {
		toolDef := skills.ToolDefinition{
			Name:        fmt.Sprintf("mcp_%s_invoke", serverName),
			Description: fmt.Sprintf("Invoke tools on the registered MCP server named %q", serverName),
			Parameters: skills.ToolParameters{
				Type: "object",
				Properties: map[string]skills.ParameterProperty{
					"tool_name": {
						Type:        "string",
						Description: "Target tool name on the MCP server",
					},
					"arguments": {
						Type:        "string",
						Description: "JSON stringified arguments to pass to the MCP tool",
					},
				},
				Required: []string{"tool_name"},
			},
			Category: "MCP Server",
		}
		skill := skills.CreateSkillFromTool(toolDef)
		skill.ID = fmt.Sprintf("mcp_server_%s", serverName)
		skill.Triggers = append(skill.Triggers, serverName)
		skill.Keywords = append(skill.Keywords, serverName)
		reg.Register(skill)
		count++
	}

	return count, nil
}

// OpenAPISpec represents a basic OpenAPI 3.0 / Swagger schema
type OpenAPISpec struct {
	OpenAPI string                  `json:"openapi"`
	Swagger string                  `json:"swagger"`
	Info    OpenAPIInfo             `json:"info"`
	Paths   map[string]OpenAPIPath  `json:"paths"`
}

// OpenAPIInfo represents metadata about the API
type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// OpenAPIPath maps HTTP methods to operations
type OpenAPIPath map[string]OpenAPIOperation

// OpenAPIOperation represents a single endpoint action
type OpenAPIOperation struct {
	Summary     string         `json:"summary"`
	Description string         `json:"description"`
	OperationID string         `json:"operationId"`
	Parameters  []any          `json:"parameters,omitempty"`
	RequestBody map[string]any `json:"requestBody,omitempty"`
}

// LoadOpenAPISpec parses an OpenAPI 3.0 or Swagger 2.0 file and registers each endpoint as a Skill
func LoadOpenAPISpec(filePath string, reg *skills.SkillRegistry) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return 0, fmt.Errorf("failed to parse OpenAPI spec JSON: %w", err)
	}

	count := 0
	for path, methods := range spec.Paths {
		for method, op := range methods {
			methodUpper := strings.ToUpper(method)
			if methodUpper != "GET" && methodUpper != "POST" && methodUpper != "PUT" && methodUpper != "DELETE" && methodUpper != "PATCH" {
				continue
			}

			opID := op.OperationID
			if opID == "" {
				cleanPath := strings.ReplaceAll(path, "/", "_")
				cleanPath = strings.ReplaceAll(cleanPath, "{", "")
				cleanPath = strings.ReplaceAll(cleanPath, "}", "")
				cleanPath = strings.Trim(cleanPath, "_")
				opID = fmt.Sprintf("%s_%s", strings.ToLower(method), cleanPath)
			}

			desc := op.Summary
			if desc == "" {
				desc = op.Description
			}
			if desc == "" {
				desc = fmt.Sprintf("HTTP %s endpoint for %s", methodUpper, path)
			}

			toolDef := skills.ToolDefinition{
				Name:        opID,
				Description: fmt.Sprintf("[%s %s] %s", methodUpper, path, desc),
				Parameters: skills.ToolParameters{
					Type: "object",
					Properties: map[string]skills.ParameterProperty{
						"path_params": {
							Type:        "string",
							Description: "JSON object containing path parameters",
						},
						"query_params": {
							Type:        "string",
							Description: "JSON object containing query string parameters",
						},
						"body": {
							Type:        "string",
							Description: "JSON payload body for the request",
						},
					},
				},
				Category: "OpenAPI Endpoint",
			}

			skill := skills.CreateSkillFromTool(toolDef)
			skill.ID = fmt.Sprintf("api_%s", strings.ToLower(opID))
			skill.Triggers = append(skill.Triggers, strings.ToLower(methodUpper), strings.ToLower(path))
			skill.Keywords = append(skill.Keywords, strings.ToLower(path))
			reg.Register(skill)
			count++
		}
	}

	return count, nil
}
