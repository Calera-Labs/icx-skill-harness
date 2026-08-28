# ICX Skill Harness

[![License](https://img.shields.io/badge/License-Apache_2.0-green.svg)](LICENSE)

Go harness for [Calera Labs Infinite Context (ICX)](https://icx.caleralabs.com). Skill libraries live in ICX. Each turn the router injects a micro-viewport of tools into your model.

## What we have

```
ICX  (A₄ Infinite Context)
 ├─ ICX Skill Harness     this repo — CLI / SDK / localhost OpenAI-compatible gateway
 │                         Skills stay in ICX. The model sees 1–2 tool schemas, not the buffet.
 └─ ICX MCP               https://icx.caleralabs.com/mcp · https://github.com/Calera-Labs/icx-mcp
                           MCP **server** for MCP **hosts** (Cursor, Claude Desktop, Windsurf, Antigravity).
```

`export-mcp` dumps `tools/list` JSON. It does not start Claude Desktop or ship a desktop connector.

Bring your own key. Get an ICX key at [dashboard.caleralabs.com](https://dashboard.caleralabs.com).

```
skills on disk  ──ingest──►  ICX space
                                  │
user prompt ──► recall + route ──► 1–2 tool schemas ──► your model
                                  │
                     optional ingest of turn state
```

## Install

Go 1.22+.

```bash
cp .env.example .env
# GEMINI_API_KEY or OPENAI_API_KEY
# ICX_API_KEY and ICX_SPACE_ID from dashboard.caleralabs.com
```

```bash
go test ./...
go build -o icx-harness ./cmd/main.go
```

Windows: `go build -o icx-harness.exe ./cmd/main.go`

## Usage

```bash
./icx-harness -cmd chat -prompt "Search arXiv for transformer papers"
./icx-harness -cmd sync -skills-dir ./skills
./icx-harness -cmd serve -port 8080
./icx-harness -cmd diagnostic -offline
```

### Gateway

`serve` listens on `127.0.0.1` and uses a gateway bearer token (printed on first run, or reused from `.icx-gateway-token`).

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://127.0.0.1:8080/v1",
    api_key="icxgw_your_gateway_token",
)
resp = client.chat.completions.create(
    model="gemini-2.5-flash",
    messages=[{"role": "user", "content": "Search arXiv for transformer papers"}],
)
print(resp.choices[0].message)
```

## Skills

Bundled files in `./skills` are schema cards. Default tool execution is a sandbox mock — it does not call vendor APIs. Load live tools with `-mcp-config` or `-openapi`.

Markdown files with YAML frontmatter in `./skills`:

```markdown
---
name: "Sandbox Echo"
id: "echo_sandbox"
category: "Starter"
description: "Echoes arguments."
triggers: ["echo", "sandbox"]
keywords: ["echo"]
---
```

## Configuration

| Variable | Purpose |
|----------|---------|
| `ICX_API_KEY` | [Calera ICX](https://dashboard.caleralabs.com) key. Omit for process-local store only. |
| `ICX_SPACE_ID` | ICX space (yours, not a shared demo) |
| `ICX_BASE_URL` | Default `https://icx.api.caleralabs.com/v1` |
| `ICX_TIMEOUT_SECONDS` | Hosted ICX HTTP timeout (default 30) |
| `BYOK_PROVIDER` | `gemini`, `openai`, `anthropic`, `deepseek`, `groq`, `ollama` |
| `GEMINI_API_KEY` / `OPENAI_API_KEY` | Model key |

Do not pass API keys on the command line. `-icx-key`, `-byok-key`, and `-gemini-key` are rejected.

Without `ICX_API_KEY`, `-cmd sync` writes to a process-local map and prints `Stored locally (not ICX)`. A hosted HTTP failure is an error, not a fake lattice success.

Pass `-crystallize` to ingest large tool outputs back into ICX.

## Commands

| `-cmd` | |
|--------|--|
| `chat` | Single prompt |
| `repl` | Interactive session |
| `serve` | Local OpenAI-compatible gateway |
| `sync` | Ingest `./skills` into ICX |
| `list` | Print the loaded registry |
| `export-mcp` | Export `tools/list` JSON |
| `diagnostic` | Routing evals |
| `benchmark` | Single-turn scorecard |
| `e2e` | Multi-turn pipelines |

## Architecture

The router scores skills by name, triggers, and keywords, boosts with ICX recall, and injects at most `-max-tools` schemas (default 2). Unmatched prompts and temporal/epistemic traps return `SAFE_REFUSAL`.

Offline diagnostic on the 115-skill eval catalog: **100.0%** overall (**100.0%** on safe refusal & always-call resistance), **97.97%** fewer prompt tokens than dump-all, and **0.1 ms** zero-allocation JIT micro-viewport routing. Scorecard: [ARCHITECTURE.md](ARCHITECTURE.md#evals). Full write-up: [CALERA_ICX_SKILL_HARNESS_PAPER.md](CALERA_ICX_SKILL_HARNESS_PAPER.md).

## Resources

- [Paper](CALERA_ICX_SKILL_HARNESS_PAPER.md)
- [Architecture](ARCHITECTURE.md)
- [ICX](https://icx.caleralabs.com)
- [Dashboard](https://dashboard.caleralabs.com)
- [ICX MCP](https://github.com/Calera-Labs/icx-mcp)
- [Calera Labs](https://caleralabs.com)

## License

Apache-2.0. See [LICENSE](LICENSE). Built by [Calera Labs](https://caleralabs.com).
