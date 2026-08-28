# Architecture

Short engineer note. The full public write-up — chef analogy, routing walkthrough, sandbox traces, glossary — is [CALERA_ICX_SKILL_HARNESS_PAPER.md](CALERA_ICX_SKILL_HARNESS_PAPER.md).

How the ICX skill harness keeps tool libraries out of the model context window.

This is the **Skill Harness** (CLI / SDK / localhost gateway). The MCP **server** MCP hosts attach is [icx-mcp](https://github.com/Calera-Labs/icx-mcp). Cursor and Claude Desktop are MCP **hosts**, not Calera desktop products.

Skill files live on disk and in a [Calera ICX](https://icx.caleralabs.com) space. On each turn the harness recalls from ICX, scores locally, and places one or two tool schemas in the prompt — a micro-viewport — then calls your model.

```
  Disk skills (markdown in ./skills)
        |
        v
  Local registry  ----------------ingest---------------->  ICX space
        |                                                     |
        |                                                recall
        v                                                     |
  Router: lexical score + ICX matches ------------------------+
        |
        v
  Micro-viewport (0–2 tool schemas) + user prompt
        |
        v
  BYOK model  -->  executor
        |
        v
  Optional ingest of large outputs back to ICX
```

## Skill bloat

Agent frameworks usually paste every tool schema into every turn. Five tools is cheap. A hundred is tens of thousands of tokens of manuals before the user's question.

On the 107-skill eval catalog, dumping every schema is about **27,224 tokens** of tool text per turn. The harness viewport is about **196 tokens** of schema. Offline diagnostic turns average **411 prompt tokens** including the user message. The unused schemas stay in ICX.

That cost shows up as token bills, time-to-first-token, and attention collisions between similar tools (`git` vs `github`). When nothing in the library matches, a dumped catalog still tends to call something. The harness returns `SAFE_REFUSAL` instead.

## Routing

1. Tokenize the prompt. Score skills by name, tool name, triggers, and keywords.
2. If `ICX_API_KEY` is set, recall against the space and boost matching documents.
3. Below threshold: `SAFE_REFUSAL`.
4. Otherwise inject at most *N* tools (default 2).

```
                            [ User prompt ]
                                  |
                   +--------------+--------------+
                   |                             |
         [ Local lexical match ]        [ ICX recall ]
                   |                             |
                   +--------------+--------------+
                                  |
                        [ Combined score ]
                                  |
                     Strong enough match?
                                  |
                 +----------------+----------------+
                 | YES                             | NO
                 v                                 v
      [ Build micro-viewport ]           [ SAFE_REFUSAL ]
      1–2 tool schemas                   no invented tool
```

Multi-step work slides the viewport: turn 1 gets skill A, turn 2 skill B. After the last tool, the harness can clear tools and let the model write.

## ICX

| Call | Role |
|------|------|
| **Ingest** | Store skill text or turn state in a space |
| **Recall** | Rank stored records for this prompt |
| **Spaces** | Partition libraries and state |

Model API keys stay on this machine. Crystallization (`-crystallize`) is off by default; when on, large tool outputs are ingested into ICX and replaced in the transcript with a pointer.

With no `ICX_API_KEY`, ingest and recall use a **process-local map** and label the result `LOCAL_FALLBACK`. That is not a hosted ICX write. With a key, HTTP errors are returned to the caller — a failed ingest is not reported as success.

Offline (`-offline`) and CI force that local path and skip hosted recall.

## BYOK

The harness does not sell model tokens. Point it at Gemini, an OpenAI-compatible endpoint, DeepSeek, Groq, or Ollama.

```bash
export GEMINI_API_KEY="your_key"
export ICX_API_KEY="your_icx_key"
export ICX_SPACE_ID="your_space"
```

```bash
./icx-harness -cmd chat -provider gemini -prompt "Search arXiv for transformer papers"
./icx-harness -cmd chat -provider ollama -byok-url "http://127.0.0.1:11434/v1" -model llama3 -prompt "..."
```

Keys belong in the environment, not on the command line.

MCP server configs load as skill metadata (`-mcp-config`). The registry exports as `tools/list` JSON (`-cmd export-mcp`). The MCP **server** MCP hosts attach is [icx-mcp](https://github.com/Calera-Labs/icx-mcp) (`https://icx.caleralabs.com/mcp`). Cursor, Claude Desktop, Windsurf, and similar apps are MCP **hosts**, not Calera desktop products. This harness is the CLI / SDK / localhost-proxy path.

## Gateway

`-cmd serve` starts an OpenAI-compatible proxy on `127.0.0.1` with full Server-Sent Events (SSE) streaming support (`stream: true`). Authenticate with the gateway token, not the model key.

The process holds the model key. CORS is limited to localhost. Dynamic endpoints include `GET /v1/models`, `GET /v1/skills` (introspection), `GET /v1/stats`, and `POST /v1/skills/register` (when `--enable-hot-register` is enabled).

## Evals

Offline snapshot from `./icx-harness -cmd diagnostic -offline` (21 cases, 115-skill evaluation catalog). Token and $/1k figures are schema-load estimates, not a billed model invoice.

| Dimension | Harness | Status |
|-----------|---------|--------|
| Distractor routing | **100%** | 🟢 Optimal |
| Always-call resistance | **100%** | 🟢 Optimal |
| Safe refusal (missing premise / temporal / unregistered traps) | **100%** | 🟢 Optimal |
| Argument accuracy | **100%** | 🟢 Optimal |
| Fault recovery | **100%** | 🟢 Optimal |
| Result grounding | **100%** | 🟢 Optimal |
| **Overall Composite Pass Rate** | **100.0%** | 🟢 **100% Grounded** |

| | Viewport | Dump-all |
|--|----------|----------|
| Prompt tokens / turn | 405.8 | 19,888 |
| Token reduction | **97.97%** | — |
| Average routing latency | **0.1 ms** | 12.4 ms |
| Est. cost / 1k tasks | $0.0304 | $1.4989 |

| Registry scale | Tools | Harness pass | Dump-all pass | Viewport tokens | Dump-all tokens | Savings |
|----------------|------:|-------------:|--------------:|----------------:|----------------:|--------:|
| Minimal | 10 | 100% | 96.8% | 212 | 2,665 | 92.1% |
| Standard | 25 | 100% | 95.0% | 218 | 6,483 | 96.6% |
| Moderate | 50 | 100% | 92.0% | 223 | 12,845 | 98.3% |
| Full catalog | 114 | 100% | 84.3% | 229 | 29,133 | 99.2% |

Dump-all token count grows with catalog size. Viewport size stays on the order of a couple of schemas. Reproduce with the command above.

## SDK

```go
registry := skills.NewSkillRegistry()
registry.LoadFromDirectory("./skills")

icxClient := icx.NewClient(icx.Config{
    BaseURL: os.Getenv("ICX_BASE_URL"),
    APIKey:  os.Getenv("ICX_API_KEY"),
    SpaceID: os.Getenv("ICX_SPACE_ID"),
})
llm := byok.NewProvider("gemini", os.Getenv("GEMINI_API_KEY"), os.Getenv("GEMINI_MODEL"), "")
jit := router.NewLatticeSkillRouter(registry, icxClient, router.DefaultRouterConfig())
runner := agent.NewRunner(registry, jit, llm, icxClient, os.Getenv("ICX_SPACE_ID"))
```
