---
name: "BigQuery SQL Optimizer & Slot Tuner"
id: "bigquery_sql_optimizer"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Optimizes BigQuery SQL queries, partition pruning, clustering keys, approximate aggregations, and slot consumption"
triggers: ["bigquery sql","partition pruning","clustering","slot reservation","bytes billed","bigquery performance"]
keywords: ["bigquery","sql","tuning","partitioning","slots","gcp"]
---

# BigQuery SQL Optimizer & Slot Tuner

> Eval fixture. Sandbox mock — not a live vendor API.

Optimizes BigQuery SQL queries, partition pruning, clustering keys, approximate aggregations, and slot consumption

```json
{
  "name": "bigquery_sql_tune",
  "description": "Tune BigQuery SQL queries for minimum slot-hours and bytes scanned",
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
