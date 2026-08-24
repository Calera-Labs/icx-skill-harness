---
name: "Qdrant Vector Database Engine"
id: "qdrant_vector_store"
execution: "sandbox-mock"
category: "Search"
description: "Stores dense high-dimensional vectors and runs payload filtered searches"
triggers: ["qdrant","vector database","payload filter","hnsw search"]
keywords: ["qdrant","vector","hnsw","payload","search"]
---

# Qdrant Vector Database Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Stores dense high-dimensional vectors and runs payload filtered searches

```json
{
  "name": "qdrant_search_points",
  "description": "Search vectors with payload filters in Qdrant",
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
