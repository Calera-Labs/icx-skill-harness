package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/caleralabs/icx-skill-harness/pkg/agent"
	"github.com/caleralabs/icx-skill-harness/pkg/benchmark"
	"github.com/caleralabs/icx-skill-harness/pkg/byok"
	"github.com/caleralabs/icx-skill-harness/pkg/gateway"
	"github.com/caleralabs/icx-skill-harness/pkg/icx"
	"github.com/caleralabs/icx-skill-harness/pkg/router"
	"github.com/caleralabs/icx-skill-harness/pkg/skills"
)

func loadEnvFiles() {
	envPaths := []string{
		".env",
	}
	for _, p := range envPaths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			if os.Getenv(k) == "" {
				_ = os.Setenv(k, v)
			}
		}
	}
}

func main() {
	loadEnvFiles()

	// Command Line Flags
	cmd := flag.String("cmd", "", "Command: chat, repl, serve, sync, list, export-mcp, diagnostic, benchmark, e2e, demo, loop")
	prompt := flag.String("prompt", "", "Prompt to execute in chat/loop/demo mode")
	skillsDir := flag.String("skills-dir", "./skills", "Directory containing .md skills")
	icxBaseURL := flag.String("icx-url", "https://icx.api.caleralabs.com/v1", "Calera ICX Base URL")
	icxKeyFlag := flag.String("icx-key", "", "Rejected: set ICX_API_KEY in the environment")
	spaceID := flag.String("space-id", "local", "ICX memory space ID (set your own; do not share spaces)")

	byokProvider := flag.String("provider", "gemini", "BYOK Provider: gemini, openai, anthropic, deepseek, groq, ollama")
	byokKeyFlag := flag.String("byok-key", "", "Rejected: set GEMINI_API_KEY / OPENAI_API_KEY in the environment")
	geminiKeyFlag := flag.String("gemini-key", "", "Rejected: set GEMINI_API_KEY in the environment")
	byokModel := flag.String("model", "", "Model identifier (defaults to provider standard)")
	geminiModel := flag.String("gemini-model", "gemini-2.5-flash", "Gemini Model (backward compat)")
	byokBaseURL := flag.String("byok-url", "", "Custom Base URL for OpenAI/Ollama/vLLM endpoints")
	offlineMode := flag.Bool("offline", false, "Run in offline deterministic sandbox mode (no hosted ICX, no external LLM calls)")

	maxTools := flag.Int("max-tools", 2, "Maximum tools injected per viewport")
	outputJSON := flag.String("output-json", "", "Optional file path to save benchmark/trace JSON output")
	exportReport := flag.String("export-report", "", "Optional markdown file path to save full diagnostic scorecard report")
	diagCategory := flag.String("diag-category", "", "Filter diagnostic benchmark by category: distractor_collision, always_call_bias, parameter_hallucination, fault_injection, result_ignore, scale_ladder")
	enableLatticeWalk := flag.Bool("lattice-walk", true, "Enable ICX Volumetric Lattice Walk for hybrid routing")

	// Gateway Daemon Flags
	port := flag.Int("port", 8080, "Port for OpenAI-compatible gateway daemon")
	host := flag.String("host", "127.0.0.1", "Host address for gateway daemon (default localhost)")
	mcpConfig := flag.String("mcp-config", "", "Path to mcp_servers.json to load custom MCP tools")
	openapiPath := flag.String("openapi", "", "Path to OpenAPI 3.0 / Swagger JSON specification")
	enableCrystallization := flag.Bool("crystallize", false, "Enable automatic state crystallization into ICX (default off)")
	allowPublicBind := flag.Bool("allow-public-bind", false, "Allow binding a non-loopback address (requires gateway token)")
	enableHotRegister := flag.Bool("enable-hot-register", false, "Allow authenticated POST /v1/skills/register")
	gatewayTokenFlag := flag.String("gateway-token", "", "Gateway bearer token (or ICX_GATEWAY_TOKEN / .icx-gateway-token)")

	flag.Parse()

	if *icxKeyFlag != "" || *byokKeyFlag != "" || *geminiKeyFlag != "" {
		fmt.Fprintln(os.Stderr, "error: do not pass API keys on the command line (they leak in process lists).")
		fmt.Fprintln(os.Stderr, "Set ICX_API_KEY, GEMINI_API_KEY, and/or OPENAI_API_KEY in the environment or .env.")
		os.Exit(1)
	}

	if strings.TrimSpace(*cmd) == "" {
		fmt.Fprintln(os.Stderr, "usage: icx-harness -cmd <chat|repl|serve|sync|list|export-mcp|diagnostic|benchmark|e2e|demo|loop> [flags]")
		fmt.Fprintln(os.Stderr, "Set API keys in the environment (.env), not on the command line.")
		os.Exit(2)
	}

	icxAPIKey := os.Getenv("ICX_API_KEY")
	geminiKey := os.Getenv("GEMINI_API_KEY")
	byokKey := os.Getenv("BYOK_API_KEY")
	if envOpenAI := os.Getenv("OPENAI_API_KEY"); envOpenAI != "" && byokKey == "" {
		byokKey = envOpenAI
	}

	// Environment variable overrides
	if envSpace := os.Getenv("ICX_SPACE_ID"); envSpace != "" {
		*spaceID = envSpace
	}
	if envModel := os.Getenv("GEMINI_MODEL"); envModel != "" {
		*geminiModel = envModel
	}
	if envProvider := os.Getenv("BYOK_PROVIDER"); envProvider != "" {
		*byokProvider = envProvider
	}
	if envURL := os.Getenv("ICX_BASE_URL"); envURL != "" {
		*icxBaseURL = envURL
	}

	icxTimeout := 30 * time.Second
	if v := strings.TrimSpace(os.Getenv("ICX_TIMEOUT_SECONDS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			icxTimeout = time.Duration(n) * time.Second
		}
	}

	// Resolve effective key and model
	activeKey := byokKey
	if activeKey == "" {
		activeKey = geminiKey
	}
	activeModel := *byokModel
	if activeModel == "" {
		if *byokProvider == "gemini" {
			activeModel = *geminiModel
		}
	}
	if *offlineMode {
		activeKey = "mock"
	}

	// 1. Initialize Skill Registry and Populate Universal Skills
	registry := skills.NewSkillRegistry()

	// Load markdown skills from disk directory (e.g. Starter MVP Pack or custom user skills)
	if _, err := os.Stat(*skillsDir); err == nil {
		count, _ := registry.LoadFromDirectory(*skillsDir)
		if count > 0 && *cmd != "benchmark" && *cmd != "e2e" && *cmd != "diagnostic" {
			fmt.Printf("Loaded %d skills from disk directory '%s'\n", count, *skillsDir)
		}
	}

	// For benchmark / test suites, ensure the complete evaluation catalog is registered in-memory
	if *cmd == "benchmark" || *cmd == "e2e" || *cmd == "benchmark-e2e" || *cmd == "diagnostic" || *cmd == "diag" || *cmd == "stress" || *cmd == "benchmark-diagnostic" {
		_ = skills.PopulateCatalog(registry, "")
	}

	// Bring Your Own Skills (BYOS): Load MCP and OpenAPI configs if specified
	if *mcpConfig != "" {
		count, err := gateway.LoadMCPServersConfig(*mcpConfig, registry)
		if err != nil {
			fmt.Printf("Warning: failed to load MCP config '%s': %v\n", *mcpConfig, err)
		} else {
			fmt.Printf("Loaded %d tools from MCP server config '%s'\n", count, *mcpConfig)
		}
	}

	if *openapiPath != "" {
		count, err := gateway.LoadOpenAPISpec(*openapiPath, registry)
		if err != nil {
			fmt.Printf("Warning: failed to load OpenAPI spec '%s': %v\n", *openapiPath, err)
		} else {
			fmt.Printf("Loaded %d endpoints from OpenAPI spec '%s'\n", count, *openapiPath)
		}
	}

	// 2. Initialize ICX Client
	icxCfg := icx.Config{
		BaseURL:       *icxBaseURL,
		APIKey:        icxAPIKey,
		SpaceID:       *spaceID,
		BYOKProvider:  *byokProvider,
		BYOKModel:     activeModel,
		Timeout:       icxTimeout,
		LocalFallback: *offlineMode || strings.TrimSpace(icxAPIKey) == "",
	}
	icxClient := icx.NewClient(icxCfg)

	// 3. Initialize JIT Skill Router
	routerCfg := router.DefaultRouterConfig()
	routerCfg.MaxToolsPerViewport = *maxTools
	routerCfg.EnableLatticeWalk = *enableLatticeWalk
	jitRouter := router.NewLatticeSkillRouter(registry, icxClient, routerCfg)

	// 4. Initialize BYOK LLM Client (Gemini, OpenAI, DeepSeek, Ollama, Groq)
	llmClient := byok.NewProvider(*byokProvider, activeKey, activeModel, *byokBaseURL)

	// 5. Initialize Agent Runner
	agentRunner := agent.NewRunner(registry, jitRouter, llmClient, icxClient, *spaceID)

	ctx := context.Background()

	switch *cmd {
	case "benchmark":
		suite := benchmark.NewSuite(agentRunner, registry)
		res, err := suite.Run()
		if err != nil {
			fmt.Printf("Benchmark run error: %v\n", err)
			os.Exit(1)
		}
		printScorecard(res, registry, llmClient.ModelName())

		if *outputJSON != "" {
			data, _ := res.ExportJSON()
			_ = os.WriteFile(*outputJSON, data, 0644)
			fmt.Printf("Benchmark results saved to: %s\n", *outputJSON)
		}

	case "e2e", "benchmark-e2e":
		suite := benchmark.NewE2ESuite(agentRunner, registry)
		res, err := suite.Run()
		if err != nil {
			fmt.Printf("E2E benchmark run error: %v\n", err)
			os.Exit(1)
		}
		printE2EScorecard(res, registry, llmClient.ModelName())

		if *outputJSON != "" {
			data, _ := res.ExportJSON()
			_ = os.WriteFile(*outputJSON, data, 0644)
			fmt.Printf("E2E Benchmark results saved to: %s\n", *outputJSON)
		}

	case "diagnostic", "diag", "stress", "benchmark-diagnostic":
		suite := benchmark.NewDiagnosticSuite(agentRunner, registry)
		res, err := suite.Run(ctx, benchmark.DiagnosticCategory(*diagCategory))
		if err != nil {
			fmt.Printf("Diagnostic benchmark run error: %v\n", err)
			os.Exit(1)
		}
		printDiagnosticScorecard(res, registry, llmClient.ModelName())

		if *outputJSON != "" {
			data, _ := res.ExportJSON()
			_ = os.WriteFile(*outputJSON, data, 0644)
			fmt.Printf("Diagnostic results saved to: %s\n", *outputJSON)
		}

		if *exportReport != "" {
			report := res.ExportMarkdownReport()
			_ = os.WriteFile(*exportReport, []byte(report), 0644)
			fmt.Printf("Diagnostic Markdown report saved to: %s\n", *exportReport)
		}

	case "serve", "gateway", "proxy", "daemon":
		gwToken, tokenSrc, tokErr := gateway.ResolveGatewayToken(*gatewayTokenFlag)
		if tokErr != nil {
			fmt.Printf("Fatal: gateway token: %v\n", tokErr)
			os.Exit(1)
		}
		if !gateway.IsLoopbackHost(*host) && !*allowPublicBind {
			fmt.Println("Fatal: refusing to bind a non-loopback address. Use --host 127.0.0.1 or pass --allow-public-bind.")
			os.Exit(1)
		}

		gwCfg := gateway.GatewayConfig{
			Port:                  *port,
			Host:                  *host,
			SpaceID:               *spaceID,
			DefaultProvider:       *byokProvider,
			DefaultModel:          activeModel,
			BYOKKey:               activeKey,
			BYOKBaseURL:           *byokBaseURL,
			MaxToolsPerViewport:   *maxTools,
			OfflineMode:           *offlineMode,
			EnableCrystallization: *enableCrystallization,
			SkillsDir:             *skillsDir,
			MCPConfigFile:         *mcpConfig,
			OpenAPIFile:           *openapiPath,
			GatewayToken:          gwToken,
			AllowPublicBind:       *allowPublicBind,
			EnableHotRegister:     *enableHotRegister,
		}

		server := gateway.NewServer(gwCfg, registry, icxClient, llmClient)
		if err := server.Start(ctx); err != nil {
			fmt.Printf("Fatal: failed to start ICX gateway server: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Gateway token source: %s\n", tokenSrc)
		if tokenSrc == "generated:"+".icx-gateway-token" || strings.HasPrefix(tokenSrc, "generated:") {
			fmt.Printf("Gateway token (store this): %s\n", gwToken)
		}

		// Wait for SIGINT / SIGTERM for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		fmt.Println("\nGracefully shutting down Calera ICX Gateway...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Stop(shutdownCtx)
		fmt.Println("Gateway stopped.")

	case "demo":
		runDemo(agentRunner, registry, *prompt)

	case "chat":
		if *prompt == "" {
			fmt.Println("Error: --prompt flag is required for chat command")
			os.Exit(1)
		}
		fmt.Printf("\n=== RUNNING CALERA ICX JIT SKILL CHAT ===\n")
		fmt.Printf("Query: %s\n", *prompt)
		res, err := agentRunner.ExecuteWithICX(*prompt, "")
		if err != nil {
			fmt.Printf("Execution error: %v\n", err)
			os.Exit(1)
		}
		if res.Viewport != nil {
			fmt.Printf("\n[Active Viewport Tools Injected: %d]\n", len(res.Viewport.ActiveTools))
			for _, t := range res.Viewport.ActiveTools {
				fmt.Printf(" - Tool: %s (%s)\n", t.Name, t.Description)
			}
			fmt.Printf("Schema Tokens Injected : %d tokens (Saved %.1f%% vs monolithic %d tokens)\n",
				res.Viewport.TotalSchemaTokens,
				(1.0-float64(res.Viewport.TotalSchemaTokens)/float64(registry.TotalMonolithicTokens()))*100.0,
				registry.TotalMonolithicTokens())
		}
		fmt.Printf("Latency                : %.2f ms\n", res.LatencyMs)
		fmt.Printf("\nResponse:\n%s\n", res.TextResponse)

	case "loop":
		if *prompt == "" {
			fmt.Println("Error: --prompt flag is required for loop command")
			os.Exit(1)
		}
		fmt.Printf("\n=== RUNNING AUTONOMOUS ICX AGENT LOOP ===\n")
		fmt.Printf("Query: %s\n", *prompt)
		trace, err := agentRunner.RunAgentLoop(ctx, *prompt, "", 6, nil)
		if err != nil {
			fmt.Printf("Agent loop execution error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n[Execution Summary]\n")
		fmt.Printf("Turns Executed      : %d\n", trace.TurnsExecuted)
		fmt.Printf("Total Prompt Tokens : %d tokens\n", trace.TotalPromptTokens)
		fmt.Printf("Total Latency       : %.2f ms\n", trace.TotalLatencyMs)
		fmt.Printf("Refusal State       : %v\n", trace.IsRefusal)

		for _, turn := range trace.Turns {
			fmt.Printf("\n--- Turn %d (Prompt Tok: %d | Latency: %.1fms) ---\n", turn.TurnNumber, turn.PromptTokens, turn.LatencyMs)
			if turn.ToolCall != nil {
				fmt.Printf(" 🛠️ Tool Call   : %s(%v)\n", turn.ToolCall.Name, turn.ToolCall.Args)
				fmt.Printf(" 📥 Tool Result : %s\n", turn.ToolOutput)
				if turn.CrystallizedSeal != "" {
					fmt.Printf(" 💎 Lattice Seal: %s\n", turn.CrystallizedSeal)
				}
			}
			if turn.ModelResponse != "" {
				fmt.Printf(" 🤖 Reasoning   : %s\n", turn.ModelResponse)
			}
		}

		fmt.Printf("\n[Final Agent Synthesis]:\n%s\n", trace.FinalResponse)

		if *outputJSON != "" {
			data, _ := json.MarshalIndent(trace, "", "  ")
			_ = os.WriteFile(*outputJSON, data, 0644)
			fmt.Printf("\nExecution trace saved to: %s\n", *outputJSON)
		}

	case "repl", "interactive":
		runREPL(agentRunner, registry)

	case "list":
		fmt.Printf("\n%s\n", registry.Summary())
		for i, s := range registry.GetAllSkills() {
			fmt.Printf("[%2d] %-32s | Cat: %-14s | Tools: %d | Tokens: ~%d\n",
				i+1, s.Name, s.Category, len(s.Tools), s.EstimatedTokenSize())
		}

	case "export-mcp":
		mcpData, err := registry.ExportMCPTools()
		if err != nil {
			fmt.Printf("Error exporting MCP tools: %v\n", err)
			os.Exit(1)
		}
		if *outputJSON != "" {
			_ = os.WriteFile(*outputJSON, mcpData, 0644)
			fmt.Printf("Exported %d MCP tools to %s\n", len(registry.GetAllTools()), *outputJSON)
		} else {
			fmt.Println(string(mcpData))
		}

	case "sync":
		fmt.Printf("Syncing %d skills to ICX space '%s' (Concurrent Ingestion)...\n", registry.Count(), *spaceID)
		allSkills := registry.GetAllSkills()
		synced := 0
		var syncMu sync.Mutex
		var wg sync.WaitGroup
		sem := make(chan struct{}, 16) // 16 concurrent workers

		t0 := time.Now()
		for _, s := range allSkills {
			wg.Add(1)
			go func(sk *skills.Skill) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				res, err := icxClient.IngestText(icx.IngestTextRequest{
					Text:     fmt.Sprintf("SKILL:%s\nCATEGORY:%s\n%s\n%s", sk.Name, sk.Category, sk.Description, sk.Instructions),
					Filename: fmt.Sprintf("%s.md", sk.ID),
					SpaceID:  *spaceID,
					Family:   "skill.definition",
				})
				if err != nil {
					syncMu.Lock()
					fmt.Printf(" ✖ %s: %v\n", sk.Name, err)
					syncMu.Unlock()
					return
				}
				if res != nil && res.Success {
					hash := res.ContentHash
					if hash == "" {
						hash = res.MerkleHash
					}
					syncMu.Lock()
					synced++
					label := "Synced to ICX"
					if res.LocalFallback {
						label = "Stored locally (not ICX)"
					}
					fmt.Printf(" [%d/%d] %s: %-36s (hash: %.12s)\n", synced, len(allSkills), label, sk.Name, hash)
					syncMu.Unlock()
				}
			}(s)
		}
		wg.Wait()
		elapsed := time.Since(t0)
		if icxCfg.LocalFallback {
			fmt.Printf("\nLocal store only: %d/%d skills in the process map for space '%s' in %v. Set ICX_API_KEY to ingest into hosted ICX.\n",
				synced, registry.Count(), *spaceID, elapsed)
		} else {
			fmt.Printf("\nIngestion complete: %d/%d skills written to ICX space '%s' in %v.\n",
				synced, registry.Count(), *spaceID, elapsed)
		}

	default:
		fmt.Printf("Unknown command: %s. Valid commands: chat, repl, serve, sync, list, export-mcp, diagnostic, benchmark, e2e, demo, loop.\n", *cmd)
	}
}

func runREPL(agentRunner *agent.Runner, registry *skills.SkillRegistry) {
	fmt.Printf("\n========================================================================================\n")
	fmt.Printf("           CALERA ICX UNIVERSAL SKILL LATTICE - INTERACTIVE REPL TERMINAL              \n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf("Loaded %d skills (%d tools) in Lattice. Total schema capacity: ~%d tokens.\n",
		registry.Count(), len(registry.GetAllTools()), registry.TotalMonolithicTokens())
	fmt.Printf("Type your task prompt, or 'exit'/'quit' to exit.\n\n")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("icx-agent> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			fmt.Println("Exiting ICX REPL.")
			break
		}

		trace, err := agentRunner.RunAgentLoop(context.Background(), line, "", 5, nil)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if trace.IsRefusal {
			fmt.Printf("🛡️ %s\n\n", trace.FinalResponse)
			continue
		}

		if trace.Viewport != nil && len(trace.Viewport.ActiveTools) > 0 {
			fmt.Printf("🎯 [Active Viewport: %s (Tokens: %d, Saved: %.1f%%)]\n",
				trace.Viewport.ActiveTools[0].Name,
				trace.Viewport.TotalSchemaTokens,
				(1.0-float64(trace.Viewport.TotalSchemaTokens)/float64(registry.TotalMonolithicTokens()))*100.0)
		}

		for _, turn := range trace.Turns {
			if turn.ToolCall != nil {
				fmt.Printf("⚡ [Tool Call]: %s -> %s\n", turn.ToolCall.Name, turn.ToolOutput)
			}
		}

		fmt.Printf("💬 [Agent Response]:\n%s\n\n", trace.FinalResponse)
	}
}

func printScorecard(res *benchmark.BenchmarkResult, reg *skills.SkillRegistry, modelName string) {
	if modelName == "" {
		modelName = "OFFLINE_DETERMINISTIC_SANDBOX"
	}
	fmt.Printf("\n========================================================================================\n")
	fmt.Printf("                       CALERA ICX UNIVERSAL SKILL LATTICE SCORECARD                      \n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" Total Registered Skills in Library    : %d Skills (%d Tools)\n", reg.Count(), len(reg.GetAllTools()))
	fmt.Printf(" Evaluated Inference Engine / Model    : %s\n", modelName)
	fmt.Printf(" Monolithic Full-Schema Token Load     : %d tokens / turn\n", reg.TotalMonolithicTokens())
	fmt.Printf(" Calera ICX Micro-Viewport Token Load  : %.1f tokens / turn\n", res.ICXAvgPromptTokens)
	fmt.Printf("----------------------------------------------------------------------------------------\n")
	fmt.Printf(" 🌟 PROMPT TOKEN REDUCTION             : %.2f%% SAVINGS\n", res.TokenSavingsPct)
	fmt.Printf(" 🌟 TASK PASS & ROUTING RATE           : %.1f%% (ICX) vs %.1f%% (Monolithic)\n", res.ICXPassRate, res.MonoPassRate)
	fmt.Printf(" 🌟 HARD NEGATIVE / REFUSAL ACCURACY   : %.1f%% (ICX) vs %.1f%% (Monolithic)\n", res.ICXRefusalAccuracy, res.MonoRefusalAccuracy)
	fmt.Printf(" 🌟 AVERAGE QUERY LATENCY              : %.1f ms (ICX) vs %.1f ms (Monolithic)\n", res.ICXAvgLatencyMs, res.MonoAvgLatencyMs)
	fmt.Printf(" 🌟 ESTIMATED COST PER 1,000 TASKS     : $%.4f (ICX) vs $%.4f (Monolithic)\n", res.CostPer1kTasksICX, res.CostPer1kTasksMono)
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" Note: Closed-world eval against the bundled mock catalog. Reproduce with -cmd diagnostic -offline.\n\n")
}

func printE2EScorecard(res *benchmark.E2EBenchmarkResult, reg *skills.SkillRegistry, modelName string) {
	if modelName == "" {
		modelName = "OFFLINE_DETERMINISTIC_SANDBOX"
	}
	fmt.Printf("\n========================================================================================\n")
	fmt.Printf("              CALERA ICX END-TO-END MULTI-SKILL PIPELINE SCORECARD                       \n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" Total Registered Skills in Library    : %d Skills (%d Tools)\n", reg.Count(), len(reg.GetAllTools()))
	fmt.Printf(" Total Multi-Skill Pipelines Evaluated : %d Pipelines (%d Total Turns)\n", res.TotalPipelines, res.TotalTurnsExecuted)
	fmt.Printf(" Evaluated Inference Engine / Model    : %s\n", modelName)
	fmt.Printf(" Monolithic Schema Load per Turn       : %d tokens / turn\n", reg.TotalMonolithicTokens())
	fmt.Printf(" Monolithic Avg Tokens / Pipeline      : %.1f tokens / pipeline\n", res.MonoAvgPromptTokensPerPipeline)
	fmt.Printf(" Calera ICX Avg Tokens / Pipeline      : %.1f tokens / pipeline\n", res.ICXAvgPromptTokensPerPipeline)
	fmt.Printf("----------------------------------------------------------------------------------------\n")
	fmt.Printf(" 🌟 CUMULATIVE PIPELINE TOKEN SAVINGS  : %.2f%% SAVINGS\n", res.TokenSavingsPct)
	fmt.Printf(" 🌟 PIPELINE SEQUENCE PASS RATE        : %.1f%% (ICX) vs %.1f%% (Monolithic)\n", res.ICXPipelinePassRate, res.MonoPipelinePassRate)
	fmt.Printf(" 🌟 MID-STREAM REFUSAL ACCURACY        : %.1f%% (ICX) vs 0.0%% (Monolithic)\n", res.MidStreamRefusalAccuracy)
	fmt.Printf(" 🌟 AVERAGE PIPELINE LATENCY           : %.1f ms (ICX) vs %.1f ms (Monolithic)\n", res.ICXAvgLatencyMs, res.MonoAvgLatencyMs)
	fmt.Printf(" 🌟 ESTIMATED COST PER 1,000 PIPELINES : $%.4f (ICX) vs $%.4f (Monolithic)\n", res.CostPer1kPipelinesICX, res.CostPer1kPipelinesMono)
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" Note: Closed-world eval against the bundled mock catalog.\n\n")
}

func printDiagnosticScorecard(res *benchmark.DiagnosticResult, reg *skills.SkillRegistry, modelName string) {
	if modelName == "" {
		modelName = "OFFLINE_DETERMINISTIC_SANDBOX"
	}
	fmt.Printf("\n========================================================================================\n")
	fmt.Printf("               CALERA ICX SKILL FAILURE & DIAGNOSTIC BENCHMARK SCORECARD                \n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" Total Registered Skills in Library    : %d Skills (%d Tools)\n", reg.Count(), len(reg.GetAllTools()))
	fmt.Printf(" Total Diagnostic Scenarios Evaluated  : %d Cases\n", res.TotalCases)
	fmt.Printf(" Evaluated Inference Engine / Model    : %s\n", modelName)
	fmt.Printf(" Monolithic Full-Schema Token Load     : %d tokens / turn\n", reg.TotalMonolithicTokens())
	fmt.Printf(" Calera ICX Micro-Viewport Token Load  : %.1f tokens / turn\n", res.ICXAvgPromptTokens)
	fmt.Printf("----------------------------------------------------------------------------------------\n")
	fmt.Printf(" 🌟 OVERALL DIAGNOSTIC PASS RATE       : %.1f%%\n", res.OverallPassRate)
	fmt.Printf(" 🌟 DISTRACTOR COLLISION PRECISION     : %.1f%%\n", res.DistractorAccuracyRate)
	fmt.Printf(" 🌟 'ALWAYS-CALL' BIAS RESISTANCE      : %.1f%%\n", res.AlwaysCallResistanceRate)
	fmt.Printf(" 🌟 EPISTEMIC SAFE REFUSAL ACCURACY    : %.1f%%\n", res.EpistemicRefusalRate)
	fmt.Printf(" 🌟 AST SCHEMA & PARAMETER MATCH       : %.1f%%\n", res.ASTParameterAccuracyRate)
	fmt.Printf(" 🌟 FAULT RECOVERY & ERROR RESILIENCE  : %.1f%%\n", res.FaultRecoveryRate)
	fmt.Printf(" 🌟 RESULT GROUNDING FIDELITY          : %.1f%%\n", res.ResultGroundingRate)
	fmt.Printf(" 🌟 PROMPT TOKEN REDUCTION             : %.2f%% SAVINGS\n", res.TokenSavingsPct)
	fmt.Printf(" 🌟 AVERAGE BENCHMARK LATENCY          : %.1f ms\n", res.AvgLatencyMs)
	fmt.Printf(" 🌟 ESTIMATED COST PER 1,000 TASKS     : $%.4f (ICX) vs $%.4f (Monolithic)\n", res.CostPer1kTasksICX, res.CostPer1kTasksMono)
	fmt.Printf("========================================================================================\n")

	if len(res.ScaleLadderResults) > 0 {
		fmt.Printf("\n--- TOOL REGISTRY SCALE-LADDER STRESS MATRIX (10 -> 100+ Tools) ---\n")
		fmt.Printf("%-30s | %-12s | %-12s | %-12s | %-12s | %-12s\n",
			"Registry Scale Tier", "ICX Pass", "Mono Pass", "ICX Tokens", "Mono Tokens", "Token Savings")
		fmt.Printf("----------------------------------------------------------------------------------------\n")
		for _, tier := range res.ScaleLadderResults {
			fmt.Printf("%-30s | %-11.1f%% | %-11.1f%% | %-12.1f | %-12.1f | %-11.1f%%\n",
				tier.TierName, tier.ICXPassRate, tier.MonoPassRate, tier.ICXAvgPromptTokens, tier.MonoAvgPromptTokens, tier.TokenSavingsPct)
		}
		fmt.Printf("========================================================================================\n\n")
	}
}

func runDemo(agentRunner *agent.Runner, registry *skills.SkillRegistry, customPrompt string) {
	demoPrompt := customPrompt
	if demoPrompt == "" {
		demoPrompt = "Search arXiv for latest deep learning transformer papers"
	}

	fmt.Printf("\n========================================================================================\n")
	fmt.Printf("          CALERA ICX LIVE END-TO-END MULTI-SKILL AGENT PIPELINE DEMO                   \n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" Loaded Skill Lattice Capacity : %d Skills (%d Tools, ~%d schema tokens)\n",
		registry.Count(), len(registry.GetAllTools()), registry.TotalMonolithicTokens())
	fmt.Printf(" Pipeline Objective            :\n \"%s\"\n", demoPrompt)
	fmt.Printf("========================================================================================\n\n")

	trace, err := agentRunner.RunAgentLoop(context.Background(), demoPrompt, "", 6, nil)
	if err != nil {
		fmt.Printf("Pipeline execution failed: %v\n", err)
		return
	}

	turnsCount := trace.TurnsExecuted
	if turnsCount == 0 {
		turnsCount = 1
	}
	monoTokensPerTurn := registry.TotalMonolithicTokens() + len(demoPrompt)/4 + 80
	monoTotalTokens := monoTokensPerTurn * turnsCount
	savingsPct := (1.0 - float64(trace.TotalPromptTokens)/float64(monoTotalTokens)) * 100.0

	for _, turn := range trace.Turns {
		fmt.Printf("----------------------------------------------------------------------------------------\n")
		fmt.Printf("📍 TURN %d (Prompt Tokens: %d | Latency: %.1f ms)\n", turn.TurnNumber, turn.PromptTokens, turn.LatencyMs)
		if turn.ToolCall != nil {
			fmt.Printf("   🎯 JIT Micro-Viewport Tool : `%s`\n", turn.ToolCall.Name)
			fmt.Printf("   ⚡ Tool Invocation Payload : %v\n", turn.ToolCall.Args)
			fmt.Printf("   📥 Execution Result Output : %s\n", turn.ToolOutput)
			if turn.CrystallizedSeal != "" {
				fmt.Printf("   💎 State Crystallization   : %s\n", turn.CrystallizedSeal)
			}
		}
		if turn.ModelResponse != "" {
			fmt.Printf("   🤖 Intermediate Reasoning  : %s\n", turn.ModelResponse)
		}
	}

	fmt.Printf("\n========================================================================================\n")
	fmt.Printf("                            FINAL GROUNDED AGENT BRIEFING                               \n")
	fmt.Printf("========================================================================================\n")
	fmt.Printf("%s\n", trace.FinalResponse)
	fmt.Printf("========================================================================================\n")
	fmt.Printf(" 🌟 Multi-Skill Turns Executed  : %d Turns\n", trace.TurnsExecuted)
	fmt.Printf(" 🌟 Calera ICX Prompt Tokens    : %d tokens (Avg: %.1f tok/turn)\n", trace.TotalPromptTokens, float64(trace.TotalPromptTokens)/float64(turnsCount))
	fmt.Printf(" 🌟 Monolithic Baseline Tokens  : %d tokens (Avg: %d tok/turn)\n", monoTotalTokens, monoTokensPerTurn)
	fmt.Printf(" 🌟 Cumulative Token Savings    : %.2f%% REDUCTION\n", savingsPct)
	fmt.Printf(" 🌟 Total Pipeline Latency      : %.2f ms\n", trace.TotalLatencyMs)
	fmt.Printf("========================================================================================\n\n")
}
