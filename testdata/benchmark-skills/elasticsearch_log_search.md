---
name: "Elasticsearch Log Search Engine"
id: "elasticsearch_log_search"
execution: "sandbox-mock"
category: "Observability"
description: "Searches distributed cluster logs and runs Lucene aggregations"
triggers: ["elasticsearch","lucene","kibana","cluster logs","log search"]
keywords: ["elasticsearch","lucene","logs","search","kibana"]
---

# Elasticsearch Log Search Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Searches distributed cluster logs and runs Lucene aggregations

```json
{
  "name": "elasticsearch_search",
  "description": "Search Elasticsearch cluster indices",
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
