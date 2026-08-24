---
name: "BigQuery Graph & GQL Engine"
id: "bigquery_graph_gql_engine"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Creates property graphs and executes Graph Query Language (GQL) path traversals directly in BigQuery"
triggers: ["bigquery graph","gql","property graph","graph query language","graph traversal"]
keywords: ["graph","gql","bigquery","nodes","edges","path"]
---

# BigQuery Graph & GQL Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Creates property graphs and executes Graph Query Language (GQL) path traversals directly in BigQuery

```json
{
  "name": "bigquery_gql_query",
  "description": "Query property graphs and topologies in BigQuery using GQL",
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
