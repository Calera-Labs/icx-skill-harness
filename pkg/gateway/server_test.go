package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/caleralabs/icx-skill-harness/pkg/byok"
	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func setupTestGatewayServer() *Server {
	reg := skills.NewSkillRegistry()
	_ = skills.PopulateCatalog(reg, "")

	icxClient := icx.NewClient(icx.Config{
		SpaceID:       "test_gateway_space",
		LocalFallback: true,
	})

	llmClient := byok.NewGeminiClient("AIzaSy_mock", "gemini-3.5-flash-lite")

	cfg := DefaultGatewayConfig()
	cfg.SpaceID = "test_gateway_space"
	cfg.OfflineMode = true
	cfg.GatewayToken = "test-gateway-token"

	return NewServer(cfg, reg, icxClient, llmClient)
}

func TestGatewayHealthAndModels(t *testing.T) {
	srv := setupTestGatewayServer()

	// 1. Health check
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	srv.handleHealthz(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /healthz, got %d", w.Code)
	}

	var healthResp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &healthResp); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}
	if healthResp["status"] != "ok" {
		t.Errorf("expected ok status, got %v", healthResp["status"])
	}
	if _, ok := healthResp["space_id"]; ok {
		t.Errorf("healthz must not disclose space_id")
	}

	// 2. Models listing
	reqModels := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	wModels := httptest.NewRecorder()
	srv.handleModels(wModels, reqModels)

	if wModels.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /v1/models, got %d", wModels.Code)
	}

	var modelsResp OpenAIModelListResponse
	if err := json.Unmarshal(wModels.Body.Bytes(), &modelsResp); err != nil {
		t.Fatalf("failed to decode models response: %v", err)
	}
	if len(modelsResp.Data) == 0 {
		t.Errorf("expected non-empty models list")
	}
}

func TestGatewayChatCompletionsDynamicBYOS(t *testing.T) {
	srv := setupTestGatewayServer()

	// Construct a request with 10 custom tools on-the-fly
	customTools := []OpenAITool{
		{
			Type: "function",
			Function: OpenAIFunctionDef{
				Name:        "custom_oracle_sql_query",
				Description: "Execute Oracle enterprise SQL database queries and transaction commits",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{"type": "string"},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: OpenAIFunctionDef{
				Name:        "custom_jira_issue_creator",
				Description: "Create and update internal enterprise Jira tickets",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"summary": map[string]any{"type": "string"},
					},
				},
			},
		},
		{
			Type: "function",
			Function: OpenAIFunctionDef{
				Name:        "custom_slack_notification_send",
				Description: "Send corporate Slack notifications to security channels",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"channel": map[string]any{"type": "string"},
					},
				},
			},
		},
	}

	// Add 10 additional filler tools to test dynamic pruning
	for i := 1; i <= 10; i++ {
		customTools = append(customTools, OpenAITool{
			Type: "function",
			Function: OpenAIFunctionDef{
				Name:        fmt.Sprintf("custom_filler_tool_%d", i),
				Description: fmt.Sprintf("Filler utility tool %d for testing scale", i),
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"input": map[string]any{"type": "string"},
					},
				},
			},
		})
	}

	chatReq := OpenAIChatCompletionRequest{
		Model: "gemini-3.5-flash-lite",
		Messages: []OpenAIMessage{
			{
				Role:    "user",
				Content: "Execute Oracle enterprise SQL database update for transaction 9981",
			},
		},
		Tools: customTools,
	}

	reqBody, _ := json.Marshal(chatReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleChatCompletions(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /v1/chat/completions, got %d. Body: %s", w.Code, w.Body.String())
	}

	var chatResp OpenAIChatCompletionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &chatResp); err != nil {
		t.Fatalf("failed to decode chat completion response: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		t.Fatalf("expected non-empty choices")
	}

	// Verify tool was selected and pruned
	if chatResp.Choices[0].Message.ToolCalls != nil && len(chatResp.Choices[0].Message.ToolCalls) > 0 {
		called := chatResp.Choices[0].Message.ToolCalls[0].Function.Name
		if !strings.Contains(called, "oracle") && !strings.Contains(called, "sql") {
			t.Errorf("expected Oracle/SQL tool to be selected, got: %s", called)
		}
	}

	// Verify token savings headers and usage metrics
	tokensSavedHeader := w.Header().Get("X-ICX-Tokens-Saved")
	if tokensSavedHeader == "" || tokensSavedHeader == "0" {
		t.Errorf("expected positive X-ICX-Tokens-Saved header")
	}

	if chatResp.Usage.ICXTokensSaved <= 0 {
		t.Errorf("expected positive ICXTokensSaved in usage stats, got %d", chatResp.Usage.ICXTokensSaved)
	}

	if chatResp.Usage.ICXSavingsPct <= 0.0 {
		t.Errorf("expected positive ICXSavingsPct in usage stats, got %.2f%%", chatResp.Usage.ICXSavingsPct)
	}

	// 3. Verify /v1/stats tracks the saved tokens
	wStats := httptest.NewRecorder()
	reqStats := httptest.NewRequest(http.MethodGet, "/v1/stats", nil)
	srv.handleStats(wStats, reqStats)

	var stats StatsSnapshot
	_ = json.Unmarshal(wStats.Body.Bytes(), &stats)
	if stats.TotalRequests != 1 {
		t.Errorf("expected 1 total request in stats, got %d", stats.TotalRequests)
	}
	if stats.TotalTokensSaved <= 0 {
		t.Errorf("expected positive total tokens saved in stats")
	}
}

func TestGatewayAuthRequired(t *testing.T) {
	srv := setupTestGatewayServer()
	handler := srv.requireAuth(srv.handleHealth)

	w := httptest.NewRecorder()
	handler(w, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	req.Header.Set("Authorization", "Bearer test-gateway-token")
	wOK := httptest.NewRecorder()
	handler(wOK, req)
	if wOK.Code != http.StatusOK {
		t.Fatalf("expected 200 with token, got %d", wOK.Code)
	}
}

func TestGatewayHotRegisterDisabled(t *testing.T) {
	srv := setupTestGatewayServer()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/skills/register", strings.NewReader(`{"name":"x"}`))
	srv.handleSkillRegister(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 when hot register is off, got %d", w.Code)
	}
}

func TestGatewayDynamicSkillRegistration(t *testing.T) {
	srv := setupTestGatewayServer()
	srv.config.EnableHotRegister = true

	newTool := skills.ToolDefinition{
		Name:        "enterprise_payroll_calculator",
		Description: "Calculate 401k matching and tax withholdings for enterprise payroll",
		Parameters: skills.ToolParameters{
			Type: "object",
			Properties: map[string]skills.ParameterProperty{
				"employee_id": {Type: "string", Description: "Corporate Employee ID"},
				"gross_pay":   {Type: "number", Description: "Gross pay amount"},
			},
			Required: []string{"employee_id", "gross_pay"},
		},
		Category: "Enterprise HR",
	}

	reqBody, _ := json.Marshal(newTool)
	req := httptest.NewRequest(http.MethodPost, "/v1/skills/register", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	srv.handleSkillRegister(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w.Code)
	}

	var regResp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &regResp)
	if regResp["status"] != "REGISTERED" {
		t.Errorf("expected REGISTERED status, got %v", regResp["status"])
	}

	// Verify tool is now discoverable in the registry
	_, _, exists := srv.baseRegistry.GetToolByName("enterprise_payroll_calculator")
	if !exists {
		t.Errorf("expected enterprise_payroll_calculator to exist in registry")
	}
}

func TestGatewayImporters(t *testing.T) {
	reg := skills.NewSkillRegistry()

	// Test OpenAPI parsing
	openAPISample := `{
		"openapi": "3.0.0",
		"info": {"title": "Company API", "version": "1.0"},
		"paths": {
			"/v1/customers": {
				"get": {"summary": "List corporate customers", "operationId": "listCustomers"},
				"post": {"summary": "Create corporate customer", "operationId": "createCustomer"}
			}
		}
	}`

	// Write temp file
	tmpFile := t.TempDir() + "/sample_openapi.json"
	_ = os.WriteFile(tmpFile, []byte(openAPISample), 0644)

	count, err := LoadOpenAPISpec(tmpFile, reg)
	if err != nil {
		t.Fatalf("failed to load openapi spec: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 endpoints loaded, got %d", count)
	}

	// Test MCP parsing
	mcpSample := `{
		"mcpServers": {
			"postgres": {
				"command": "npx",
				"args": ["-y", "@modelcontextprotocol/server-postgres", "postgresql://localhost/db"]
			}
		}
	}`
	tmpMCP := t.TempDir() + "/sample_mcp.json"
	_ = os.WriteFile(tmpMCP, []byte(mcpSample), 0644)

	mcpCount, err := LoadMCPServersConfig(tmpMCP, reg)
	if err != nil {
		t.Fatalf("failed to load MCP config: %v", err)
	}
	if mcpCount != 1 {
		t.Errorf("expected 1 MCP server loaded, got %d", mcpCount)
	}
}

func TestGatewayChatCompletionsStreaming(t *testing.T) {
	srv := setupTestGatewayServer()

	chatReq := OpenAIChatCompletionRequest{
		Model: "gemini-3.5-flash-lite",
		Messages: []OpenAIMessage{
			{
				Role:    "user",
				Content: "Explain quantum computing briefly",
			},
		},
		Stream: true,
	}

	reqBody, _ := json.Marshal(chatReq)
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.handleChatCompletions(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from streaming completions, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		t.Errorf("expected text/event-stream content type, got: %s", contentType)
	}

	bodyStr := w.Body.String()
	if !strings.Contains(bodyStr, "data:") {
		t.Errorf("expected SSE data: chunks in body")
	}
	if !strings.Contains(bodyStr, "data: [DONE]") {
		t.Errorf("expected data: [DONE] sentinel in stream output")
	}
}

func TestGatewayListSkills(t *testing.T) {
	srv := setupTestGatewayServer()

	req := httptest.NewRequest(http.MethodGet, "/v1/skills", nil)
	w := httptest.NewRecorder()
	srv.handleListSkills(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /v1/skills, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode /v1/skills response: %v", err)
	}

	if resp["total_skills"] == nil || resp["total_skills"].(float64) <= 0 {
		t.Errorf("expected positive total_skills count")
	}
}
