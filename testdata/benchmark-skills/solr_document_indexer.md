---
name: "Apache Solr Document Indexer"
id: "solr_document_indexer"
execution: "sandbox-mock"
category: "Search"
description: "Indexes text documents and runs full-text BM25 queries"
triggers: ["solr","indexer","search engine","solr collection"]
keywords: ["solr","search","indexer","bm25","lucene"]
---

# Apache Solr Document Indexer

> Eval fixture. Sandbox mock — not a live vendor API.

Indexes text documents and runs full-text BM25 queries

```json
{
  "name": "solr_index_ops",
  "description": "Manage Solr document collections",
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
