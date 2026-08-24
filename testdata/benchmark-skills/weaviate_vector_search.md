---
name: "Weaviate Vector & Hybrid Search Engine"
id: "weaviate_vector_search"
execution: "sandbox-mock"
category: "Search"
description: "Performs ANN vector similarity searches with cosine distance filters"
triggers: ["weaviate","vector","embedding","similarity","hybrid search"]
keywords: ["weaviate","vector","embedding","hybrid","ann"]
---

# Weaviate Vector & Hybrid Search Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Performs ANN vector similarity searches with cosine distance filters

```json
{
  "name": "weaviate_vector_query",
  "description": "Perform hybrid vector queries on Weaviate",
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
