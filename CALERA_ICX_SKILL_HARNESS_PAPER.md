# Using Calera ICX to Kill Skill Bloat: An Open Harness for Just-in-Time Tool Viewports

**A layman’s deep dive into skill bloat, micro-viewports, Bring-Your-Own-Key models, and what this Go client actually does with Calera ICX.**

**What this document is:** A public write-up of the open-source **ICX skill harness**—a Go client that stores skill libraries in **Calera ICX**, injects only the tools needed for the current turn, and calls **your** model keys (BYOK).

**What this document is not:** A description of ICX internals. Calera ICX is treated here as a hosted memory and recall service: ingest, recall, and spaces. How ICX stores bytes on the server is out of scope.

**Where this sits:**

```
ICX  (A₄ Infinite Context)
 ├─ ICX Skill Harness     this paper — CLI / SDK / localhost OpenAI gateway
 └─ ICX MCP               https://icx.caleralabs.com/mcp
                           MCP server for MCP hosts (Cursor, Claude Desktop, Windsurf, Antigravity).
```

**Methods note:** Benchmark numbers below come from a **closed-world sandbox**. The bundled catalog and tool executor are mocks. They are not live Stripe, Kubernetes, Vault, or SEC APIs. Reproduce with `./icx-harness -cmd diagnostic -offline`. See the repository README.

*Calera Labs · August 2026 · Apache-2.0 harness*

---

```
  +--------------------------------------------------------------------------------------------------+
  |                         WHAT THE HARNESS DOES WITH ICX                                           |
  |                                                                                                  |
  |   Skill files on disk  ----ingest---->  Calera ICX space (your library, your partition)          |
  |                                                                                                  |
  |   User prompt  -->  ICX recall + local routing                                                   |
  |                -->  1–2 tool schemas in the prompt (a micro-viewport)                            |
  |                -->  Your model (Gemini / GPT / Claude / DeepSeek / Ollama / …)                   |
  |                -->  Optional: store turn state back in ICX so the chat does not bloat            |
  +--------------------------------------------------------------------------------------------------+
```

---

## Table of contents

1. [The problem we actually solved](#1-the-problem-we-actually-solved)
2. [The chef analogy](#2-the-chef-analogy)
3. [Skill bloat, in practical terms](#3-skill-bloat-in-practical-terms)
4. [What ICX is for (and what we will not describe)](#4-what-icx-is-for-and-what-we-will-not-describe)
5. [Architecture of the open harness](#5-architecture-of-the-open-harness)
6. [Bring-your-own-key](#6-bring-your-own-key)
7. [Sandbox benchmarks (honest)](#7-sandbox-benchmarks-honest)
8. [Worked sandbox traces](#8-worked-sandbox-traces)
9. [How to run it](#9-how-to-run-it)
10. [What we are not claiming](#10-what-we-are-not-claiming)
11. [The road ahead](#11-the-road-ahead)
12. [Glossary](#12-glossary)

---

## 1. The problem we actually solved

In the race to build agents that do real work—query a warehouse, open a pull request, search papers, file a ticket—the industry hit a boring, expensive bottleneck. We call it **skill bloat**, or the **context window tax**.

When you give an assistant tools, most frameworks treat the model’s prompt as a dumping ground. They paste the **full instruction manual** for every tool into **every turn**: names, descriptions, JSON parameter schemas, sometimes examples. Five tools is harmless. One hundred is not.

As a library grows from a handful of tools to 100 or 500 skills, the model may have to read **tens of thousands of tokens of boilerplate** before it reads the user’s question. You pay for that reread. You wait for it. And the model’s attention is spent on manuals it does not need this turn.

```
+----------------------------------------------------------------------------------------------------+
|                               THE MONOLITHIC AGENT BOTTLENECK                                      |
|                                                                                                    |
|   1. Cost:     Every turn re-sends the entire tool catalog. A five-step loop bills five catalogs.  |
|   2. Latency:  Time-to-first-token grows with unused schemas, not with the question.               |
|   3. Confusion: Overlapping names (git vs github, sql vs postgres) get worse as the dump grows.    |
|   4. Action bias: When nothing matches, bloated agents still tend to call something.               |
+----------------------------------------------------------------------------------------------------+
```

We did not invent a new model. We built a **client around Calera ICX** so the library can live **outside** the context window.

**The play, in five steps:**

1. Register skills locally (markdown with YAML frontmatter, optional MCP or OpenAPI import).
2. **Ingest** those definitions into an ICX **space** (optional for offline demos, intended for real use).
3. On each user turn, **recall** from ICX and/or score locally, then inject **one or two** tool schemas—a **micro-viewport**.
4. Call a BYOK model with that tiny payload. You keep paying the model vendor; you are not locked to one.
5. Optionally **ingest** large tool outputs back into ICX so conversation history does not recreate the same bloat.

Anyone can write another harness against the same ICX API. This repository is a reference client, not the only legal way to talk to ICX. That is the point: ICX is the memory; the harness is a USB cable you can replace.

---

## 2. The chef analogy

To see the design without drowning in schemas, picture a dinner service.

### The old way: the overburdened chef

You hire a brilliant chef. Before they can ask what you want to eat, you strap a backpack to their shoulders containing **500 heavy binders**: how to fly a 747, how to weld underwater, German grammar, tax filing, billing APIs, container orchestration, protein databases. Then you say:

> “Chef, scramble two eggs.”

Before they crack an egg, they must leaf through the backpack. Aviation diagrams. Welding tolerances. A chapter on invoices. Somewhere in the pile are three lines about butter and a pan.

What happens in that kitchen?

1. **The chef is slow.** Most of the time is spent lifting binders, not cooking. That is time-to-first-token: the model cannot start answering until it has ingested the dump.
2. **The meal is expensive.** You pay the chef for every page they were forced to scan, not just for the omelet. That is prompt-token billing on unused tool schemas.
3. **The chef makes silly mistakes.** With 50,000 pages of noise, motor oil can look a lot like butter. Two tools named something like `git` and `github` sit next to each other. The wrong one gets grabbed.
4. **Impossible orders become inventions.** You ask for *moon-dust soufflé*. There is no such recipe in the kitchen. A trustworthy chef would say so. An overwhelmed chef, staring at 500 binders that *almost* mention dust, baking, and silver, will guess—and serve you baking powder and dirt.

```
                    THE OVERBURDENED CHEF (monolithic agent)

         [ Manual: How to Fly a 747 ]      [ Manual: Deep Sea Welding ]
         [ Manual: Quantum Physics  ]      [ Manual: Tax Code         ]
         [ Manual: German Grammar   ]      [ Manual: Billing APIs     ]
                               \        /
                                v      v
                           ( @ _ @ )  <-- Chef trying to make an omelet
                           "I have 500 binders on the counter.
                            I can barely find the salt."
```

That backpack-on-every-order pattern is how most agent frameworks “support tools” today. It works at three tools. It collapses at three hundred.

### The new way: ICX as the vault, the harness as the runner

Now the same restaurant, rebuilt around Calera ICX and this client.

- The 500 binders live in a **library vault down the hall**. That vault is an **ICX space**. You do not drag it into the kitchen. You do not reprint it on every ticket. It is still *there*—searchable, durable, yours.
- Beside the chef stands a fast runner: the **JIT micro-viewport router** in this harness.
- You say “scramble two eggs.” The runner checks the vault (ICX recall) and the index cards already in the kitchen (local skill files). It pulls **one card**: eggs, butter, pan. That card is maybe 180 tokens. It goes on the counter. The other 499 binders stay in the vault.
- The chef reads ten lines, cooks, serves. No backpack.

```
                   ICX + THIS HARNESS

     [ 500 encyclopedias stored in your ICX space ]
                        |
     Ticket: "Scramble two eggs with butter"
                        |
         [ Micro-viewport router ]
                        |
        Places 1 index card on the counter: "Egg scrambling"
                        |
                        v
                   ( ^ _ ^ )  <-- Chef focused, unburdened
                   "The vault has the library. I only need this card."
```

### Moon-dust soufflé

The interesting case is not eggs. Eggs are easy. The interesting case is the order the kitchen **cannot** fill.

You ask for moon-dust soufflé. The runner looks in ICX. No card. It looks in the local registry. No trigger, no tool name, no keyword that clears the bar. It does **not** hand the chef a vaguely related pastry manual and hope. It answers at the door:

> “We do not carry moon dust. I cannot fulfill this request.”

That is **`SAFE_REFUSAL`**. It is a **harness** behavior. It happens *before* the chef (the language model) is asked to invent a recipe. It is also *not* a claim that models never hallucinate once they are called. If you put a real tool in the viewport and the model lies about the arguments, that is still a model problem. What we refuse to do is *force* a tool call when the library has no matching skill.

In the old kitchen, missing recipes were the most dangerous tickets: 500 binders create a bias toward *doing something*. In the sandbox evals, the dump-everything baseline hallucinated a tool on every trap prompt we ran. The harness refused. That is why the vault-and-runner split exists—not as theater, but as a place to say no.

### Why the analogy maps onto software

| Kitchen | Software |
|---------|----------|
| Binders | Tool JSON schemas and skill markdown |
| Counter space | The model context window (you pay for every page on it) |
| Vault down the hall | A Calera ICX **space** |
| Runner who fetches one card | This harness’s router (local score + ICX recall) |
| Single index card | The **micro-viewport** (typically 1–2 tools) |
| “We don’t carry moon dust” | `SAFE_REFUSAL` |
| Paying the chef by the page | Prompt tokens billed by the model vendor |
| The vault’s membership | Your ICX key and space—not the chef’s wages |

The chef can be Gemini one night and a local Ollama model the next. That is BYOK. The vault does not care which chef is on the line. The binders do not get reprinted when you change chefs.

---

## 3. Skill bloat, in practical terms

Models do not have hands, mice, or a native connection to your warehouse. When an agent “queries Postgres,” it emits a structured function call. Your process runs SQL and stuffs the result back into the next turn.

To make that possible, the developer writes a schema:

1. **Name** — e.g. `execute_sql_query`
2. **Description** — what it is for, and often what it is *not* for
3. **Parameters** — types, required fields, nested objects

The host copies those schemas into the hidden prompt. If you have three tools, that is a few hundred tokens of tax. If you have a “universal workstation”—twenty databases, a DevOps drawer, SaaS APIs, scientific lookups—the tax becomes the meal.

In our sandbox catalog of **107 mock skills**, dumping every schema measured about **27,224 tokens** of tool text per turn. The user’s question might be 35 tokens. The harness viewport measured about **196 tokens**. That ratio is why we built this. It is a **schema-size** measurement on fixtures, not a promise about your production OpenAPI files.

```
+----------------------------------------------------------------------------------------------------+
|                       PROMPT CONTEXT COMPARISON (ONE USER TURN, SANDBOX CATALOG)                   |
|                                                                                                    |
|  Dump-every-schema agent:                                                                          |
|  [========================================================================] 27,224 tool tokens     |
|  [==] ~35 user tokens                                                                              |
|                                                                                                    |
|  ICX harness micro-viewport:                                                                       |
|  [==] ~160 viewport tool tokens                                                                    |
|  [==] ~35 user tokens                                                                              |
|                                                                                                    |
|  The unused 27,000 tokens did not vanish. They live in ICX (or on disk). They are not in the bill. |
+----------------------------------------------------------------------------------------------------+
```

### Three penalties, then a fourth

**1. Money (token inflation).** Cloud vendors charge for prompt tokens. A five-step ticket that resends 27k schemas each turn is ~135k tokens of *static manuals* before any real work. At team volume that is not a rounding error. In the sandbox heuristic we used, dump-all vs viewport on that catalog was on the order of **$0.60 vs $0.01 per thousand single-turn tasks**—a back-of-envelope using a simple $/million-token rate, **not** an ICX invoice. ICX is a separate seat.

**2. Speed (TTFT).** Before a model emits its first character, it must ingest the input. A 27k-token dump is a different job from a 200-token viewport. Our local sandbox sometimes showed dump-all “faster” in-process because it was not a real 27k cloud round trip. In production, large prompts are slower. Do not quote the sandbox’s 96 ms dump-all number as “monolithic is quicker.”

**3. Confusion (attention).** Models do not have infinite, uniform focus. Pack a prompt with dense JSON and “lost in the middle” is not a metaphor: similar tools steal from each other. The router’s job is to *not put the distractors on the counter*.

**4. Missing tools (the moon-dust problem).** Ask for an IBM Quantum QPU when no such skill exists. A trustworthy system should refuse. In the sandbox, the harness refused; dump-all still tried to call a tool. That is a routing result on fixtures. It is also the failure mode that scared us enough to put refusal in the client instead of hoping the model would abstain.

### How agents actually call tools (one paragraph more)

The loop is always: **schema in → model proposes call → host executes → result in → model speaks.** Skill bloat lives entirely in “schema in.” Shrink that, and the rest of the loop can stay standard: OpenAI-shaped tools, MCP, LangChain, Cursor. That is why the gateway speaks `/v1/chat/completions`. We did not need a new religion for step two. We needed the manuals off the counter.

---

## 4. What ICX is for (and what we will not describe)

**Calera ICX** (Infinite Context) is a **hosted memory service**. From this harness’s point of view it has a small, boring, essential job:

| Call | What the harness uses it for |
|------|------------------------------|
| **Ingest** | Put skill text or turn state into a **space** so it survives this process |
| **Recall** | Ask ICX which stored records match this prompt, to boost routing |
| **Spaces** | Keep one team’s library and state off another team’s |

That is the entire ICX story this paper tells. We **do not** document how the hosted service stores data. If you need that, it belongs in Calera product documentation, not in an open client write-up.

Think of ICX the way you already think of “a database the agent is allowed to remember against.” You would not paste the whole database into the system prompt. You would query it. ICX is that query surface for **skills and state**.

### Product contrast (not internals)

People ask “why not a vector database?” Fair question. This is the *product* answer, not a bake-off of indexes:

- **ICX is the memory this client is built to speak.** Skills and turn state share a space. You ingest once; you recall when a prompt arrives.
- **The harness still routes locally** (names, triggers, keywords) so **offline** demos and CI still work. Unreachable ICX falls back to a process-local map. That map is not ICX. It is a spare flashlight.
- **Online, recall can boost which skill wins.** The model still only sees the viewport. ICX does not dump the library into the prompt. If it did, we would have rebuilt the backpack.

### Privacy defaults in the public client

These are footguns we already stepped on in private and then nailed shut:

- Your **model API key stays on your machine**. The open client does not send it to ICX.
- ICX ingest/recall uses an **ICX API key** and a **space id you choose**. Do not share a demo space with the internet.
- Writing every turn into ICX is **off** unless you pass `-crystallize`. Silent default ingest is how strangers’ prompts end up in one pile.

ICX is the vault. The vault has a key. The chef’s wages (BYOK) are a different wallet.

---

## 5. Architecture of the open harness

The repo is a Go CLI plus an optional localhost OpenAI-compatible gateway. We chose Go for a fast local router and a single binary, not because ICX requires Go. Python, TypeScript, or a one-file script can speak the same ingest/recall calls. We want that.

```
  Disk skills (markdown in ./skills)
        |
        v
  Local registry  ----------------ingest---------------->  ICX space
        |                                                     |
        |                                                recall (if keyed)
        v                                                     |
  Router: lexical score + ICX matches ------------------------+
        |
        v
  Micro-viewport (0–2 tool schemas) + user prompt
        |
        v
  BYOK model  -->  sandbox mock executor, or your real HTTP/MCP executor
        |
        v
  Optional ingest of large outputs back to ICX
```

### 5.1 Storage is not attention

Traditional agents treat the context window as both warehouse and workbench. Our split is:

- **Storage (capacity):** as many skill files as you want, on disk and/or in ICX.
- **Attention (this turn):** only the viewport.

ICX is how storage stays durable across processes, laptops, and nights. The local registry is how a clone of this repo still routes when you have not set an ICX key yet.

### 5.2 The micro-viewport

The router returns a small result: matched skills, the active tool JSON, a token estimate, a confidence score, and either a viewport or `SAFE_REFUSAL`.

If nothing in the library matches strongly enough, the harness **does not** call the model with a decoy tool. It returns a refusal. That is the epistemic behavior we wanted in the *client*, because hoping the model will abstain after you dumped 107 manuals on it did not work in our traps.

A viewport is allowed to be empty on a synthesis turn: tools already ran; now the model just writes the briefing.

### 5.3 How routing uses ICX

1. Tokenize the prompt. Score skills by **name, tool name, triggers, keywords**.
2. If an ICX key is configured, **recall** against the space. When ICX returns a matching skill document, add a score boost.
3. If the best score is below a threshold, refuse.
4. Otherwise inject at most *N* tools (default 2).

**Strong anchors** (exact names and triggers) matter more than a vague keyword overlap. That is how `git status` does not win a Kubernetes prompt just because both mention “cluster” in passing—or whatever accidental collision your catalog invents.

Offline (`-offline`) and missing `ICX_API_KEY` skip hosted recall. Ingest in that mode writes a process-local map labeled `LOCAL_FALLBACK`, not a lattice success. If a key is set and ICX is unreachable, ingest/recall return an error. CI uses the local path. Online, ICX is the difference between “index cards in this process” and “the library you ingested last week, still there.”

```
                       ROUTING DECISION (PUBLIC DESCRIPTION)

                            [ User prompt ]
                                  |
                   +--------------+--------------+
                   |                             |
         [ Local lexical match ]        [ ICX recall ]
         (names, triggers, tools)       (if API key set)
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
      * 1–2 tool schemas                 * No fake tool call
      * Typically < 200 tokens           * Model need not invent a skill
```

### 5.4 Multi-step prompts

Real work is rarely one tool. A four-step prompt should not load four manuals at once.

The runner **slides the viewport**: turn 1 gets skill A, turn 2 skill B, and so on. After the last tool, it can clear tools and let the model summarize.

```
+----------------------------------------------------------------------------------------------------+
|               SLIDING VIEWPORT (IDEA, NOT A LIVE VENDOR DEMO)                                      |
|                                                                                                    |
|  Turn 1: [ Skill A only ]  --> ran                                                                 |
|  Turn 2:                  [ Skill B only ]  --> ran                                                |
|  Turn 3:                                       [ Skill C only ]  --> ran                           |
|  Turn 4:  [ no tools — write the briefing ]                                                        |
+----------------------------------------------------------------------------------------------------+
```

In the **sandbox**, A/B/C return canned JSON tagged `[SANDBOX MOCK TOOL]`. In your fork, replace the executor with real HTTP or MCP. The viewport logic does not care, as long as you do not paste the whole catalog again.

### 5.5 Optional state ingest (“crystallization”)

A tool that returns a 10,000-word blob will recreate bloat if you leave that blob in chat forever. If `-crystallize` is on, the harness **ingests** the blob into ICX and leaves a short pointer in the transcript. The model can keep the facts it already used; later turns do not pay for the raw dump again.

Default in the public binary is **off**. Silent ingest into a shared space is how you mix tenants. Opt in when you have your own space and you know what you are storing.

### 5.6 Gateway (local only)

`-cmd serve` starts an OpenAI-compatible proxy on **127.0.0.1**. It requires a **gateway token**—a local password, not your model key and not your ICX key. Point LangChain or a Python `OpenAI()` client at `http://127.0.0.1:8080/v1` and pass that gateway token as `api_key`.

The daemon **holds your model key**. Do not bind it to the public internet unless you understand that anyone who can reach it can spend that key. CORS is limited to localhost origins. Hot skill registration over HTTP is off unless you explicitly enable it.

---

## 6. Bring-your-own-key

The harness does not sell model tokens. You point it at Gemini, an OpenAI-compatible endpoint, DeepSeek, Groq, or Ollama. Switch chefs; keep the vault.

```
+==================================================================================================+
|                               BYOK: ONE ROUTER, MANY MODELS                                      |
|                                                                                                  |
|                             +------------------------------+                                     |
|                             |  ICX space + this harness    |                                     |
|                             +------------------------------+                                     |
|                                            |                                                     |
|                      +---------------------+---------------------+                               |
|                      v                     v                     v                               |
|               [ Google Gemini ]     [ OpenAI-compatible ]  [ Local Ollama ]                      |
+==================================================================================================+
```

Set keys in the environment (`.env` is gitignored). **Do not put secrets on the command line**; they land in process lists and shell history.

```bash
export GEMINI_API_KEY="your_key"
export ICX_API_KEY="your_icx_key"          # optional; omit for local-only routing
export ICX_SPACE_ID="your_private_space"   # your space, not a shared demo
```

```bash
./icx-harness -cmd chat -provider gemini -prompt "Search arXiv for transformer papers"
./icx-harness -cmd chat -provider ollama -byok-url "http://127.0.0.1:11434/v1" -model llama3 -prompt "..."
```

**MCP:** you can import an MCP server *config* as skill metadata (we do not copy command lines into tool descriptions the model will read) and export the registry as `tools/list` JSON. That does not make this repo a marketplace. The official **ICX MCP** server (`https://icx.caleralabs.com/mcp`, [icx-mcp](https://github.com/Calera-Labs/icx-mcp)) is a separate Calera project. You attach it to an MCP **host** such as Cursor, Claude Desktop, Windsurf, or Antigravity. Claude Desktop is Anthropic's app — not a Calera desktop connector. This harness is the CLI / SDK / localhost-proxy path for people who already have an OpenAI client.

---

## 7. Sandbox benchmarks (honest)

We compared two strategies on the **same** eval catalog:

- **Monolithic / dump-all:** inject all 107 mock schemas every turn.
- **Harness:** inject the routed viewport (and use ICX recall when a key is configured).

**Catalog:** in-memory fixtures (`PopulateCatalog` / `testdata/benchmark-skills`). **Executor:** pattern-matched JSON, always prefixed `[SANDBOX MOCK TOOL]`.

A later 21-case offline diagnostic (including the 10→107 tool scale ladder) is in [ARCHITECTURE.md](ARCHITECTURE.md#evals). That diagnostic is the **canonical published scorecard**. Reproduce with `./icx-harness -cmd diagnostic -offline`.

### 7.1 Diagnostic scorecard (21 cases, 107 mock skills)

Closed-world sandbox. Token and $/1k figures are schema-load estimates, not a billed model invoice.

| Dimension | Harness |
|-----------|---------|
| Distractor routing | 100% |
| Always-call resistance | 80% |
| Safe refusal (missing premise / unregistered skill) | **66.7%** |
| Argument accuracy | 100% |
| Fault recovery | 100% |
| Result grounding | 100% |
| **Overall** | **95.2%** |

| | Viewport | Dump-all |
|--|----------|----------|
| Prompt tokens / turn | 410.6 | 27,321 |
| Token reduction | 98.50% | — |
| Est. cost / 1k tasks | $0.0308 | $2.0491 |

Safe refusal is the weak dimension: two of three missing-premise / unregistered-skill traps are caught. Do not quote 100% routing from an older 19-prompt suite as the public number.

These pass rates mean: **on this fixture set, the router picked the labeled tool (or refused traps).** They do not mean field accuracy on live APIs. They do not mean the model never hallucinates *after* a tool is selected.

Illustrative tasks in that suite (all mocks): a 10-K-shaped lookup, a git-diff-shaped call, SQL-shaped, container-shaped, a billing-shaped name, IAM-shaped, plus **traps** (quantum QPU, Solana contract) that must refuse.

### 7.2 Multi-turn pipelines (6 sandbox scripts, 28 turns)

Same catalog. Pipelines are **scripted sequences** of mock tools (finance-shaped, DevOps-shaped, literature-shaped, data-shaped) plus a **mid-stream trap**.

Cumulative schema-token savings in that run: **~96.7%** (about 128k tokens of dumped schemas per pipeline vs about 4k with viewports). Sequence match on those scripts: harness 100% vs dump-all 83.3% in that run. Mid-stream trap: harness refused the missing skill; dump-all did not.

Again: sequence match on **named mocks**, not “we healed production Kubernetes.”

### 7.3 Economics, scoped correctly

If you currently pay to resend ~27k tool tokens per turn, shrinking that to ~200 tokens is a large **prompt** saving. Dollar figures use a simple $/million-token heuristic. They are **not** a Calera bill. ICX is a hosted seat (flat rate with published quotas), separate from the model vendor.

The harness does not take a cut of “tokens saved.” It is an open client. The vault can be billed; the USB cable is free.

A back-of-envelope from the sandbox catalog, for intuition only:

- Dump-all, five turns: on the order of 5 × (prompt + 27k schemas).
- Viewport, five turns: 5 × (prompt + ~200 schemas), plus whatever you later store in ICX if you opt into crystallization.

Scale that by your real vendor price and your real schema sizes. Do not put `$1,112,520 / year` on a slide unless you have measured *your* traffic.

### 7.4 Scaling the library

As mock library size went from 10 to 100+ tools in the diagnostic ladder, dump-all token count grew with N. Viewport size stayed on the order of a couple of schemas. That is the point: **the prompt should not scale with unused skills**. ICX (or disk) can scale. The counter must not.

---

## 8. Worked sandbox traces

The following traces are **from the mock executor**. Treat tickers, accessions, protein scores, and invoice ids as **fixtures**, not as audited market or clinical data.

### 8.1 Sliding viewport on a four-step script

**Prompt (sandbox):** extract a 10-K-shaped metric, run a DCF-shaped calc, “update” SQL, “send” Slack.

The harness injected **one** mock tool per turn, then a synthesis turn with no tools:

```
Turn 1  viewport: sec_edgar_query          -> [SANDBOX MOCK TOOL] { ticker, margin, ... }
Turn 2  viewport: valuation_dcf_calc       -> [SANDBOX MOCK TOOL] { dcf fixture }
Turn 3  viewport: postgres_executor        -> [SANDBOX MOCK TOOL] { committed fixture }
Turn 4  viewport: slack_send_message       -> [SANDBOX MOCK TOOL] { sent fixture }
Turn 5  viewport: (none)                   -> model writes a briefing from the fixtures
```

Schema tokens stayed in the low hundreds per turn instead of 27,224. That is the sliding-card version of the chef analogy: four different index cards, never the whole backpack.

### 8.2 Moon dust in the middle of a real-looking ticket

**Prompt (sandbox):** two in-catalog mock steps, then “run quantum annealing on an IBM QPU.”

Turns 1–2 routed to mocks (billing-shaped, cache-shaped). Turn 3 had **no matching skill** in the registry or in ICX. The harness returned `SAFE_REFUSAL` and did not invent a quantum tool.

That is the moon-dust soufflé, served as JSON instead of dirt. It is the behavior we want when **both** the local library and the ICX space lack a capability.

---

## 9. How to run it

```bash
cp .env.example .env
# Fill GEMINI_API_KEY or OPENAI_API_KEY. Optional: ICX_API_KEY and ICX_SPACE_ID.

go test ./...
go build -o icx-harness ./cmd/main.go

# Closed-world eval (no model, no ICX required)
./icx-harness -cmd diagnostic -offline

# Starter pack on disk (8 honest schemas, still mock-executed in this repo)
./icx-harness -cmd chat -prompt "Search arXiv for transformer papers"

# Local gateway (127.0.0.1 + gateway token printed or stored in .icx-gateway-token)
./icx-harness -cmd serve -port 8080

# Push local skill text into your ICX space (needs ICX_API_KEY)
./icx-harness -cmd sync -skills-dir ./skills
```

Windows: `go build -o icx-harness.exe ./cmd/main.go` then `.\icx-harness.exe ...`.

**Adding a skill:** drop a markdown file with YAML frontmatter in `./skills` (see `skills/echo_sandbox.md`). Do not add trademarked “official integration” names unless you have a real executor and the right to use the name. The 100+ files under `testdata/benchmark-skills` are **eval mocks**, not a partner catalog.

Example frontmatter shape:

```markdown
---
name: "Sandbox Echo"
id: "echo_sandbox"
category: "Starter"
execution: "sandbox-mock"
description: "Labeled sandbox skill that echoes arguments."
triggers: ["echo", "sandbox"]
keywords: ["echo", "mock"]
---
```

**Go SDK sketch** (keys from env, never from source):

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

Python against the local gateway:

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://127.0.0.1:8080/v1",
    api_key="icxgw_your_gateway_token",  # gateway token, not the model key
)
resp = client.chat.completions.create(
    model="gemini-2.5-flash",
    messages=[{"role": "user", "content": "Search arXiv for transformer papers"}],
)
print(resp.choices[0].message)
```

---

## 10. What we are not claiming

- We are not publishing how ICX stores data on the server. This paper only uses ingest, recall, and spaces.
- We are not shipping 107 live enterprise integrations. `/skills` is a small starter pack; the large set is **eval mocks**.
- Sandbox routing numbers are **closed-world fixture scores**. The published diagnostic is **95.2% overall**, with safe refusal at **66.7%**. That is not “zero hallucination in the wild.”
- Ingest receipts are **content hashes**, useful as ids—not a full audit-product spec.
- The localhost gateway is **not** a multi-tenant cloud.
- Dollar savings in the scorecard are **heuristics**, not ICX pricing.

What we **are** claiming: you can keep a large skill library, use **ICX as the memory**, show the model only a viewport, keep **BYOK** keys local, and refuse when nothing in the library matches. Anyone else can wire the same ICX API. This repo is one kitchen. The vault is ICX.

---

## 11. The road ahead

Three things this architecture is for, without promising a science-fiction calendar:

1. **Libraries that outgrow the prompt.** Hundreds or thousands of skills should make ICX busier, not the context window fatter.
2. **Refusal as a feature.** High-stakes tools need a client that can say “not in the library” without spending a model call on a guess.
3. **Provider independence.** Skills and state live in ICX and on disk. The chef is replaceable. That is BYOK.

Compatible third-party harnesses are welcome. Speak ingest, recall, spaces, and viewports. Do not send model keys to ICX. Bind local by default.

---

## 12. Glossary

- **BYOK:** Bring your own model API key. You pay the model vendor.
- **Calera ICX:** Hosted infinite-context memory: ingest and recall against a space. This paper uses only that API-level description.
- **Context window:** The model’s short-term working memory. Everything in it is billed and attended to.
- **Crystallization (optional):** Ingesting a turn’s output into ICX so the chat does not keep the full blob.
- **Function calling / tool use:** The model emits a structured call; the host runs it.
- **Harness:** This open Go client (parser, registry, router, gateway, BYOK adapters).
- **Micro-viewport:** The 0–2 tool schemas actually placed in the prompt this turn. The index card.
- **MCP:** Model Context Protocol. ICX also has a dedicated MCP server for IDEs; this harness is the CLI/SDK path.
- **Monolithic / dump-all agent:** An architecture that pastes every tool schema into every turn. The backpack kitchen.
- **SAFE_REFUSAL:** Harness-level refusal when no registered skill matches. Happens before a needless model tool call. Moon-dust at the door.
- **Skill bloat / context window tax:** Paying attention and money to unused tool schemas every turn.
- **Space:** ICX partition for one library and its state. Use your own. Do not share it.
- **Time-to-first-token (TTFT):** Delay from sending a prompt to the first output character.
- **Token:** Billing/attention unit of model text (often ~4 characters in English).

---

*Open harness: Apache-2.0. Hosted ICX: Calera Labs. Model tokens: your vendor. Copyright © 2026 Calera Labs.*
