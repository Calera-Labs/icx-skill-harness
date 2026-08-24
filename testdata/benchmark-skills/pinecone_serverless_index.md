---
name: "Pinecone Serverless Vector Indexer"
id: "pinecone_serverless_index"
execution: "sandbox-mock"
category: "Search"
description: "Manages Pinecone serverless vector namespaces and executes sub-50ms top-K embeddings retrieval"
triggers: ["pinecone","pinecone index","vector namespace","top-k vectors"]
keywords: ["pinecone","vector","serverless","embeddings","topk"]
---

# Pinecone Serverless Vector Indexer

> Eval fixture. Sandbox mock — not a live vendor API.

Manages Pinecone serverless vector namespaces and executes sub-50ms top-K embeddings retrieval

```json
{
  "name": "pinecone_query_vectors",
  "description": "Query vectors and metadata in Pinecone serverless indices",
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
