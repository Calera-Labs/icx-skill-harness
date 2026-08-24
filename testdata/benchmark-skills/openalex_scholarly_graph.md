---
name: "OpenAlex Global Scholarly Graph"
id: "openalex_scholarly_graph"
execution: "sandbox-mock"
category: "Research & Literature"
description: "Navigates 250M+ scholarly works, author h-index bibliometrics, institutional citations, and Open Access DOIs"
triggers: ["openalex","scholarly graph","author h-index","citation count","doi lookup","bibliometrics"]
keywords: ["openalex","scholar","citations","hindex","doi","publications"]
---

# OpenAlex Global Scholarly Graph

> Eval fixture. Sandbox mock — not a live vendor API.

Navigates 250M+ scholarly works, author h-index bibliometrics, institutional citations, and Open Access DOIs

```json
{
  "name": "openalex_query_graph",
  "description": "Query the OpenAlex scholarly knowledge graph and bibliometrics",
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
