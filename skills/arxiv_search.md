---
name: "ArXiv Search"
id: "arxiv_search"
category: "Starter"
execution: "sandbox-mock"
description: "Sandbox mock of an arXiv metadata search. This harness does not call export.arxiv.org."
triggers: ["arxiv", "preprint", "research paper", "arxiv search"]
keywords: ["arxiv", "paper", "preprint", "research", "mock"]
---

# ArXiv Search

Default tool execution in this repo is a sandbox mock. Wire a live HTTP executor if you want real arXiv results.

```json
{
  "name": "arxiv_search_papers",
  "description": "Sandbox mock: search arXiv preprint metadata by query string",
  "parameters": {
    "type": "object",
    "properties": {
      "query": {"type": "string", "description": "Search query"}
    },
    "required": ["query"]
  }
}
```
