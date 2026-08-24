---
name: "Apollo Federation Subgraph Router"
id: "graphql_apollo_federation"
execution: "sandbox-mock"
category: "API"
description: "Composes subgraphs into a unified supergraph gateway schema"
triggers: ["apollo","federation","supergraph","subgraph","apollo router"]
keywords: ["apollo","federation","supergraph","subgraph","graphql"]
---

# Apollo Federation Subgraph Router

> Eval fixture. Sandbox mock — not a live vendor API.

Composes subgraphs into a unified supergraph gateway schema

```json
{
  "name": "apollo_compose_subgraph",
  "description": "Validate and compose Apollo Federation subgraphs",
  "parameters": {
    "type": "object",
    "properties": {
      "query": {"type": "string", "description": "Primary action query"},
      "options": {"type": "string", "description": "Optional parameters"}
    },
    "required": ["query"]
  }
}
```
