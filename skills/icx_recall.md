---
name: "ICX Recall"
id: "icx_recall"
category: "Starter"
execution: "sandbox-mock"
description: "Sandbox mock of an ICX recall skill. The router may call hosted ICX separately when ICX_API_KEY is set."
triggers: ["recall", "icx search", "what did we store"]
keywords: ["icx", "recall", "memory", "search", "mock"]
---

# ICX Recall

The default tool executor returns a fixture. Hosted recall is the ICX client used by the router, not this mock tool.

```json
{
  "name": "icx_recall",
  "description": "Sandbox mock of querying an ICX space",
  "parameters": {
    "type": "object",
    "properties": {
      "query": {"type": "string", "description": "Natural language or keyword query"},
      "top_k": {"type": "integer", "description": "Maximum matches to return"}
    },
    "required": ["query"]
  }
}
```
