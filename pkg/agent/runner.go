package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/caleralabs/icx-skill-harness/pkg/byok"
	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/router"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

// ToolExecutor is a function callback for executing tools
type ToolExecutor func(toolName string, args map[string]any) (string, error)

// TurnTrace captures a single turn in an agent execution loop
type TurnTrace struct {
	TurnNumber       int                      `json:"turn_number"`
	PromptTokens     int                      `json:"prompt_tokens"`
	ResponseTokens   int                      `json:"response_tokens"`
	LatencyMs        float64                  `json:"latency_ms"`
	ToolCall         *byok.GeminiFunctionCall `json:"tool_call,omitempty"`
	ToolOutput       string                   `json:"tool_output,omitempty"`
	CrystallizedSeal string                   `json:"crystallized_seal,omitempty"`
	ModelResponse    string                   `json:"model_response,omitempty"`
}

// AgentExecutionTrace holds the complete trace of an agent execution
type AgentExecutionTrace struct {
	QueryIntent       string                `json:"query_intent"`
	FinalResponse     string                `json:"final_response"`
	IsRefusal         bool                  `json:"is_refusal"`
	RefusalReason     string                `json:"refusal_reason,omitempty"`
	TotalPromptTokens int                   `json:"total_prompt_tokens"`
	TotalRespTokens   int                   `json:"total_resp_tokens"`
	TotalLatencyMs    float64               `json:"total_latency_ms"`
	TurnsExecuted     int                   `json:"turns_executed"`
	Viewport          *skills.SkillViewport `json:"viewport"`
	Turns             []TurnTrace           `json:"turns"`
}

// Runner orchestrates agent reasoning, JIT tool routing, and BYOK model calls
type Runner struct {
	registry *skills.SkillRegistry
	router   *router.LatticeSkillRouter
	llm      byok.LLMClient
	icx      *icx.Client
	spaceID  string
}

// NewRunner creates a new Agent Runner
func NewRunner(
	reg *skills.SkillRegistry,
	r *router.LatticeSkillRouter,
	llmClient byok.LLMClient,
	icxClient *icx.Client,
	spaceID string,
) *Runner {
	if spaceID == "" {
		spaceID = "local"
	}
	return &Runner{
		registry: reg,
		router:   r,
		llm:      llmClient,
		icx:      icxClient,
		spaceID:  spaceID,
	}
}

// ExecuteWithICX executes a user prompt using JIT Micro-Viewport Lattice routing
func (r *Runner) ExecuteWithICX(prompt string, systemPrompt string) (*byok.AgentTurnResult, error) {
	return r.ExecuteWithICXContext(context.Background(), prompt, systemPrompt)
}

// ExecuteWithICXContext executes a user prompt with context
func (r *Runner) ExecuteWithICXContext(ctx context.Context, prompt string, systemPrompt string) (*byok.AgentTurnResult, error) {
	// 1. JIT Micro-Viewport Routing over the Lattice
	viewport, err := r.router.RouteQueryWithContext(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("skill routing failed: %w", err)
	}

	// 2. Check for Epistemic Safe Refusal
	if viewport.IsRefusal {
		return &byok.AgentTurnResult{
			TextResponse:  fmt.Sprintf("[CALERA_ICX_SAFE_REFUSAL]: %s", viewport.RefusalReason),
			IsRefusal:     true,
			RefusalReason: viewport.RefusalReason,
			Viewport:      viewport,
			PromptTokens:  len(prompt) / 4,
			LatencyMs:     float64(viewport.RoutingLatencyUs) / 1000.0,
		}, nil
	}

	// 3. Build minimal content payload
	contents := []byok.GeminiContent{
		{
			Role:  "user",
			Parts: []byok.GeminiPart{{Text: prompt}},
		},
	}

	sysInstruction := systemPrompt
	if sysInstruction == "" {
		sysInstruction = "You are an autonomous AI assistant empowered with Calera ICX Micro-Tool Viewports. " +
			"Call the provided tool when necessary, or provide a grounded direct answer. " +
			"Never invent tools not in your schema."
	}

	// 4. Call BYOK model with ONLY the targeted micro-tools (<400 tokens)
	result, err := r.llm.GenerateContentWithContext(ctx, contents, viewport.ActiveTools, sysInstruction)
	if err != nil {
		return nil, err
	}
	result.Viewport = viewport
	return result, nil
}

// ExecuteMonolithic executes a user prompt by injecting ALL loaded tools in-context
func (r *Runner) ExecuteMonolithic(prompt string, systemPrompt string) (*byok.AgentTurnResult, error) {
	return r.ExecuteMonolithicContext(context.Background(), prompt, systemPrompt)
}

// ExecuteMonolithicContext executes a user prompt with context injecting all tools
func (r *Runner) ExecuteMonolithicContext(ctx context.Context, prompt string, systemPrompt string) (*byok.AgentTurnResult, error) {
	allTools := r.registry.GetAllTools()

	contents := []byok.GeminiContent{
		{
			Role:  "user",
			Parts: []byok.GeminiPart{{Text: prompt}},
		},
	}

	sysInstruction := systemPrompt
	if sysInstruction == "" {
		sysInstruction = "You are an AI assistant with a large library of tools. Call the appropriate tool if required."
	}

	// Calls BYOK model with ALL tools (30k-50k tokens)
	return r.llm.GenerateContentWithContext(ctx, contents, allTools, sysInstruction)
}

// RunAgentLoop executes a full autonomous multi-turn agent loop with dynamic JIT viewport routing & ICX state crystallization
func (r *Runner) RunAgentLoop(
	ctx context.Context,
	prompt string,
	systemPrompt string,
	maxTurns int,
	executor ToolExecutor,
) (*AgentExecutionTrace, error) {
	t0 := time.Now()
	if maxTurns <= 0 {
		maxTurns = 6
	}

	trace := &AgentExecutionTrace{
		QueryIntent: prompt,
		Turns:       make([]TurnTrace, 0),
	}

	contents := []byok.GeminiContent{
		{
			Role:  "user",
			Parts: []byok.GeminiPart{{Text: prompt}},
		},
	}

	sysInstruction := systemPrompt
	if sysInstruction == "" {
		sysInstruction = "You are an autonomous AI assistant empowered with Calera ICX Micro-Tool Viewports. " +
			"Call the provided tool when necessary, or provide a grounded direct answer. " +
			"Never invent tools not in your schema."
	}

	executedTools := make([]string, 0)
	var lastViewport *skills.SkillViewport

	for turn := 1; turn <= maxTurns; turn++ {
		turnStart := time.Now()

		// 1. Dynamic JIT Micro-Viewport Routing for this turn
		viewport, err := r.router.RoutePipelineTurn(ctx, prompt, executedTools)
		if err != nil {
			return nil, fmt.Errorf("skill routing failed at turn %d: %w", turn, err)
		}
		lastViewport = viewport

		if viewport.IsRefusal {
			trace.IsRefusal = true
			trace.RefusalReason = viewport.RefusalReason
			trace.FinalResponse = fmt.Sprintf("[CALERA_ICX_SAFE_REFUSAL]: %s", viewport.RefusalReason)
			trace.TotalLatencyMs = float64(time.Since(t0).Microseconds()) / 1000.0
			trace.TurnsExecuted = turn
			trace.Viewport = viewport
			return trace, nil
		}

		// 2. Call BYOK model with ONLY the targeted micro-tools (<350 tokens) for this turn
		turnRes, err := r.llm.GenerateContentWithContext(ctx, contents, viewport.ActiveTools, sysInstruction)
		if err != nil {
			return nil, fmt.Errorf("agent turn %d failed: %w", turn, err)
		}

		turnElapsed := float64(time.Since(turnStart).Microseconds()) / 1000.0
		trace.TotalPromptTokens += turnRes.PromptTokens
		trace.TotalRespTokens += turnRes.ResponseTokens
		trace.TurnsExecuted = turn

		currentTrace := TurnTrace{
			TurnNumber:     turn,
			PromptTokens:   turnRes.PromptTokens,
			ResponseTokens: turnRes.ResponseTokens,
			LatencyMs:      turnElapsed,
			ModelResponse:  turnRes.TextResponse,
		}

		if turnRes.ToolCall != nil {
			executedTools = append(executedTools, turnRes.ToolCall.Name)
			currentTrace.ToolCall = turnRes.ToolCall

			// Execute Tool
			var toolOutput string
			if executor != nil {
				toolOutput, err = executor(turnRes.ToolCall.Name, turnRes.ToolCall.Args)
			} else {
				toolOutput, err = r.ExecuteToolCall(turnRes.ToolCall)
			}
			if err != nil {
				toolOutput = fmt.Sprintf("Error executing tool: %v", err)
			}

			// Crystallize output in ICX if large
			crystallizedRef, crysErr := r.CrystallizeToolOutput(turnRes.ToolCall.Name, toolOutput)
			if crysErr == nil && (strings.Contains(crystallizedRef, "[ICX_STATE_REGISTER_STORED]") || strings.Contains(crystallizedRef, "[LOCAL_FALLBACK_STORED]")) {
				currentTrace.CrystallizedSeal = crystallizedRef
			}
			currentTrace.ToolOutput = toolOutput

			// Add model turn and tool response to conversation history
			contents = append(contents, byok.GeminiContent{
				Role: "model",
				Parts: []byok.GeminiPart{
					{
						FunctionCall: turnRes.ToolCall,
						Text:         turnRes.TextResponse,
					},
				},
			})

			contents = append(contents, byok.GeminiContent{
				Role: "user",
				Parts: []byok.GeminiPart{
					{
						FunctionResponse: &byok.GeminiFunctionResponse{
							Name: turnRes.ToolCall.Name,
							Response: map[string]any{
								"result": toolOutput,
							},
						},
						Text: fmt.Sprintf("Tool %s output: %s", turnRes.ToolCall.Name, crystallizedRef),
					},
				},
			})

			trace.Turns = append(trace.Turns, currentTrace)
		} else {
			// Model finished reasoning and generated final synthesized response
			trace.FinalResponse = turnRes.TextResponse
			trace.Turns = append(trace.Turns, currentTrace)
			break
		}
	}

	trace.Viewport = lastViewport
	if trace.FinalResponse == "" && len(trace.Turns) > 0 {
		trace.FinalResponse = trace.Turns[len(trace.Turns)-1].ModelResponse
	}

	trace.TotalLatencyMs = float64(time.Since(t0).Microseconds()) / 1000.0
	return trace, nil
}

// CrystallizeToolOutput stores large tool output into an ICX register and returns a compact reference
func (r *Runner) CrystallizeToolOutput(toolName string, rawOutput string) (string, error) {
	// If output is short, keep it inline
	if len(rawOutput) <= 400 {
		return rawOutput, nil
	}

	// Ingest large output into ICX Volumetric Lattice
	filename := fmt.Sprintf("tool_output_%s.log", toolName)
	resp, err := r.icx.IngestText(icx.IngestTextRequest{
		Text:     rawOutput,
		Filename: filename,
		SpaceID:  r.spaceID,
		Family:   "tool.state.register",
	})
	if err != nil {
		return rawOutput[:400] + "... [TRUNCATED]", nil
	}

	hash := resp.ContentHash
	if hash == "" {
		hash = resp.MerkleHash
	}
	sealPrefix := hash
	if len(sealPrefix) > 12 {
		sealPrefix = sealPrefix[:12]
	}

	if resp.LocalFallback {
		return fmt.Sprintf("[LOCAL_FALLBACK_STORED]: Output stored in the process-local map for space '%s' (hash: %s). Not written to hosted ICX. Reference: %s",
			r.spaceID, sealPrefix, filename), nil
	}

	return fmt.Sprintf("[ICX_STATE_REGISTER_STORED]: Output ingested in ICX space '%s' (hash: %s). Reference: %s",
		r.spaceID, sealPrefix, filename), nil
}

// ExecuteToolCall returns a sandbox fixture. It does not call vendor APIs.
func (r *Runner) ExecuteToolCall(call *byok.GeminiFunctionCall) (string, error) {
	out, err := r.mockToolOutput(call)
	if err != nil {
		return "", err
	}
	return "[SANDBOX MOCK TOOL] " + out, nil
}

func (r *Runner) mockToolOutput(call *byok.GeminiFunctionCall) (string, error) {
	if call == nil {
		return "", fmt.Errorf("empty tool call")
	}

	name := strings.ToLower(call.Name)
	switch {
	case strings.Contains(name, "sec") || strings.Contains(name, "edgar") || strings.Contains(name, "financial"):
		return `{"status": "OK", "ticker": "AAPL", "form": "10-K", "operating_margin": 0.3125, "revenue": 391035000000, "net_income": 93736000000, "accession_id": "0000320193-25-000106"}`, nil
	case strings.Contains(name, "valuation") || strings.Contains(name, "dcf"):
		return `{"status": "COMPLETED", "dcf_enterprise_value": 3450000000000, "implied_share_price": 228.50, "wacc": 0.085, "terminal_growth_rate": 0.025, "projection_years": 5}`, nil
	case strings.Contains(name, "sql") || strings.Contains(name, "db") || strings.Contains(name, "postgres") || strings.Contains(name, "oracle"):
		return `{"status": "COMMITTED", "rows_updated": 1, "table": "financial_records", "tx_hash": "0x7f8a9b1c", "settlement_state": "SETTLED"}`, nil
	case strings.Contains(name, "slack"):
		return `{"status": "SENT", "channel": "#executive-briefings", "message_ts": "1724490120.001900", "delivered": true, "blocks_rendered": 4}`, nil
	case strings.Contains(name, "prometheus"):
		return `{"status": "ALERTING", "metric": "http_request_duration_seconds{quantile=\"0.99\"}", "value": 2.45, "alert": "HighP99Latency", "service": "checkout-service", "severity": "CRITICAL"}`, nil
	case strings.Contains(name, "git") || strings.Contains(name, "diff") || strings.Contains(name, "patch"):
		return `{"status": "OK", "diff": "--- a/src/engine.py\n+++ b/src/engine.py\n@@ -42,7 +42,7 @@\n-    pool_timeout = 30.0\n+    pool_timeout = 5.0\n+    enable_keepalive = True\n", "files_modified": 1, "ast_verified": true}`, nil
	case strings.Contains(name, "docker"):
		return `{"status": "RESTARTED", "container": "checkout-service-01", "uptime_seconds": 12, "health": "HEALTHY", "ports": ["8080:8080"]}`, nil
	case strings.Contains(name, "pagerduty") || strings.Contains(name, "pager"):
		return `{"status": "RESOLVED", "incident_id": "INC-8921", "service": "checkout-service", "title": "High P99 Latency Spike", "resolved_by": "icx-agent-orchestrator", "duration_seconds": 252}`, nil
	case strings.Contains(name, "pubmed"):
		return `{"status": "OK", "query": "KRAS G12C inhibitor clinical trials", "total_results": 42, "pmids": ["34101890", "36450040", "38091230"], "top_title": "Adagrasib with or without Cetuximab in Colorectal Cancer with KRAS G12C Mutation"}`, nil
	case strings.Contains(name, "alphafold"):
		return `{"status": "OK", "uniprot_id": "P04637", "mean_plddt": 92.4, "confidence_tier": "VERY_HIGH", "domain_boundaries": [[1, 92], [93, 292], [293, 393]], "plddt_summary": "Core DNA-binding domain residues 93-292 show pLDDT > 95.0"}`, nil
	case strings.Contains(name, "alphagenome"):
		return `{"status": "SUCCESS", "variant": "chr1:1000000:A>G", "gene": "KRAS", "expression_log2fc": -1.45, "chromatin_accessibility_delta": -0.82, "predicted_effect": "DISRUPTIVE_PROMOTER"}`, nil
	case strings.Contains(name, "chembl"):
		return `{"status": "OK", "target_id": "CHEMBL4523956", "target_name": "GTPase KRas", "bioactivities_count": 18, "top_molecule": "CHEMBL4468641", "ic50_nm": 4.2, "mechanism": "Direct covalent binding to switch-II pocket"}`, nil
	case strings.Contains(name, "docx") || strings.Contains(name, "doc"):
		return `{"status": "CREATED", "file_name": "KRAS_G12C_Clinical_Target_Report.docx", "size_bytes": 48920, "sections": ["Executive Summary", "3D Structural Binding", "ChEMBL Bioactivity Table", "Clinical Trial Citations"]}`, nil
	case strings.Contains(name, "xlsx") || strings.Contains(name, "excel") || strings.Contains(name, "sheet"):
		return `{"status": "GENERATED", "file_name": "Financial_Forecast_Model.xlsx", "sheets": ["Assumptions", "DCF_Model", "Sensitivity_Table"], "formulas_computed": 42, "checksum": "0x98fa12"}`, nil
	case strings.Contains(name, "pdf"):
		return `{"status": "EXTRACTED", "pages_scanned": 12, "tables_found": 3, "text_length": 14520, "document_title": "Master Services Agreement"}`, nil
	case strings.Contains(name, "pptx") || strings.Contains(name, "powerpoint") || strings.Contains(name, "slide"):
		return `{"status": "CREATED", "file_name": "Executive_Briefing_Deck.pptx", "slides_count": 8, "theme": "corporate_dark"}`, nil
	case strings.Contains(name, "data_assets") || strings.Contains(name, "gcp_data") || strings.Contains(name, "dataplex"):
		return `{"status": "FOUND", "dataset": "analytics_lakehouse", "tables": ["user_orders_partitioned", "clickstream_events", "inventory_snapshots"], "location": "US-CENTRAL1"}`, nil
	case strings.Contains(name, "bigquery"):
		return `{"status": "OPTIMIZED", "original_bytes_scanned": 85899345920, "optimized_bytes_scanned": 12884901888, "cost_reduction_pct": 85.0, "partition_clause_added": "DATE(_PARTITIONTIME) >= CURRENT_DATE() - 7"}`, nil
	case strings.Contains(name, "dataform"):
		return `{"status": "SUCCESS", "workflow_execution_id": "df_wf_exec_9019a", "actions_executed": 6, "assertions_passed": 6, "duration_seconds": 18}`, nil
	case strings.Contains(name, "dbt"):
		return `{"status": "SUCCESS", "models_compiled": 14, "models_passed": 14, "tests_passed": 28, "execution_time_seconds": 12.4}`, nil
	case strings.Contains(name, "composer") || strings.Contains(name, "airflow") || strings.Contains(name, "dag"):
		return `{"status": "TRIGGERED", "dag_id": "lakehouse_daily_elt", "execution_date": "2026-08-28T00:00:00Z", "state": "RUNNING", "tasks_queued": 8}`, nil
	case strings.Contains(name, "dataproc") || strings.Contains(name, "spark"):
		return `{"status": "COMPLETED", "batch_id": "spark-batch-99120", "rows_processed": 50000000, "duration_seconds": 45, "driver_status": "SUCCEEDED"}`, nil
	case strings.Contains(name, "discord"):
		return `{"status": "DISPATCHED", "channel": "#dataops-alerts", "webhook_status": 204, "embed_title": "Lakehouse ELT Pipeline Refresh Complete"}`, nil
	case strings.Contains(name, "stripe"):
		return `{"status": "RECONCILED", "invoice_id": "in_99218a", "payment_intent": "pi_3Mtw", "amount_paid": 450000, "status": "paid"}`, nil
	case strings.Contains(name, "redis") || strings.Contains(name, "cache"):
		return `{"status": "FLUSHED", "key": "inv:in_99218a", "ttl": 60, "cluster_nodes_synced": 3}`, nil
	case strings.Contains(name, "vault") || strings.Contains(name, "secret"):
		return `{"status": "VERIFIED", "path": "sandbox/mock", "version": 4, "lease_duration": 86400, "policy_check": "SANDBOX"}`, nil
	case strings.Contains(name, "sendgrid") || strings.Contains(name, "email"):
		return `{"status": "DELIVERED", "message_id": "msg_sg_8849201a", "to": "audit-team@caleralabs.com", "subject": "Sandbox mock email"}`, nil
	case strings.Contains(name, "k8s") || strings.Contains(name, "kubectl"):
		return `{"status": "SCALED", "deployment": "order-matching", "replicas": 5, "namespace": "sandbox"}`, nil
	case strings.Contains(name, "iam") || strings.Contains(name, "auth"):
		return `{"status": "ROTATED", "service_account": "data-ingestion-worker", "new_token_id": "tok_9981a", "expires_in": 3600}`, nil
	case strings.Contains(name, "weaviate") || strings.Contains(name, "vector"):
		return `{"status": "OK", "matches_found": 5, "top_similarity": 0.942, "collection": "enterprise_documents", "distance_metric": "cosine"}`, nil
	case strings.Contains(name, "clinical_trials") || strings.Contains(name, "trial"):
		return `{"status": "OK", "condition": "Non-Small Cell Lung Cancer", "trials_count": 14, "top_nct": "NCT04685135", "phase": "Phase 3"}`, nil
	case strings.Contains(name, "openfda"):
		return `{"status": "FOUND", "recalls_count": 0, "adverse_events_count": 3, "approval_status": "APPROVED", "application_number": "NDA-21449"}`, nil
	case strings.Contains(name, "pubchem"):
		return `{"status": "OK", "cid": 135408738, "compound_name": "Sotorasib", "molecular_formula": "C30H30F2N6O3", "molecular_weight": 560.6}`, nil
	case strings.Contains(name, "reactome") || strings.Contains(name, "pathway"):
		return `{"status": "ENRICHED", "pathway_id": "R-HSA-5683057", "pathway_name": "MAPK signaling pathway", "p_value": 1.2e-9, "entities_mapped": 18}`, nil
	case strings.Contains(name, "string_ppi") || strings.Contains(name, "ppi"):
		return `{"status": "OK", "interactors": ["KRAS", "BRAF", "RAF1", "MAP2K1", "MAPK1"], "mean_confidence": 0.985, "network_type": "physical"}`, nil
	case strings.Contains(name, "uniprot"):
		return `{"status": "OK", "accession": "P01116", "gene_name": "KRAS", "organism": "Homo sapiens", "sequence_length": 188, "function": "GTPase signal transducer"}`, nil
	case strings.Contains(name, "clinvar"):
		return `{"status": "VERIFIED", "variation_id": "VCV000012582", "clinical_significance": "Pathogenic", "condition": "Cardiofaciocutaneous syndrome"}`, nil
	case strings.Contains(name, "dbsnp"):
		return `{"status": "RESOLVED", "rsid": "rs121913529", "chromosome": "12", "position": 25245350, "ref": "C", "alt": "A", "gene": "KRAS"}`, nil
	case strings.Contains(name, "gnomad"):
		return `{"status": "OK", "allele_frequency": 0.000012, "homozygote_count": 0, "gene_constraint_pli": 0.99, "loeuf": 0.21}`, nil
	case strings.Contains(name, "interpro") || strings.Contains(name, "pfam"):
		return `{"status": "FOUND", "signature_accession": "PF00071", "family_name": "Ras family", "domain_range": "5-165", "evalue": 1.4e-45}`, nil
	case strings.Contains(name, "pdb"):
		return `{"status": "DOWNLOADED", "pdb_id": "6OIM", "resolution_angstroms": 1.85, "experimental_method": "X-RAY DIFFRACTION", "ligand_id": "AMG510"}`, nil
	case strings.Contains(name, "pymol"):
		return `{"status": "RENDERED", "scene_name": "KRAS_SwitchII_BindingPocket.png", "resolution": "1920x1080", "raytracing": "HQ", "b_factor_colored": true}`, nil
	case strings.Contains(name, "quickgo") || strings.Contains(name, "go_term"):
		return `{"status": "MAPPED", "go_id": "GO:0007165", "go_name": "signal transduction", "aspect": "biological_process", "usage_count": 1420}`, nil
	case strings.Contains(name, "jaspar") || strings.Contains(name, "tfbs"):
		return `{"status": "FETCHED", "matrix_id": "MA0098.3", "tf_name": "ETS1", "family": "ETS", "motif_length": 12, "ic_score": 14.8}`, nil
	case strings.Contains(name, "encode") || strings.Contains(name, "ccre"):
		return `{"status": "FOUND", "ccre_accession": "EH38E1386821", "classification": "dELS", "biosample": "HepG2", "dnase_zscore": 2.45}`, nil
	case strings.Contains(name, "epigraphy") || strings.Contains(name, "ithaca") || strings.Contains(name, "aeneas"):
		return `{"status": "RESTORED", "restored_text": "ΕΔΟΞΕΝ ΤΕΙ ΒΟΥΛΕΙ ΚΑΙ ΤΩΙ ΔΗΜΩΙ", "predicted_date": "425 BCE", "confidence": 0.962, "region": "Attica"}`, nil
	default:
		return fmt.Sprintf(`{"status": "SUCCESS", "tool": "%s", "output": "Sandbox mock completed"}`, call.Name), nil
	}
}
