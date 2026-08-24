---
name: "Milvus Distributed Vector Database"
id: "milvus_vector_cluster"
execution: "sandbox-mock"
category: "Search"
description: "Scales billion-scale vector indexes with GPU acceleration"
triggers: ["milvus","ann search","vector index","billion scale vector"]
keywords: ["milvus","vector","gpu","ann","cluster"]
---

# Milvus Distributed Vector Database

> Eval fixture. Sandbox mock — not a live vendor API.

Scales billion-scale vector indexes with GPU acceleration

```json
{
  "name": "milvus_ann_search",
  "description": "Execute ANN vector search on Milvus cluster",
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
