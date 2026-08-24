---
name: "ChromaDB Embedded Vector Store"
id: "chromadb_local_embeddings"
execution: "sandbox-mock"
category: "Search"
description: "Stores, persists, and queries local collection embeddings and document metadata with ChromaDB"
triggers: ["chromadb","chroma","local embeddings","collection query"]
keywords: ["chroma","chromadb","embeddings","local","rag"]
---

# ChromaDB Embedded Vector Store

> Eval fixture. Sandbox mock — not a live vendor API.

Stores, persists, and queries local collection embeddings and document metadata with ChromaDB

```json
{
  "name": "chromadb_collection_search",
  "description": "Search ChromaDB embedded document collections",
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
