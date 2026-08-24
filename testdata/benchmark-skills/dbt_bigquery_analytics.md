---
name: "dbt BigQuery Analytics Engineer"
id: "dbt_bigquery_analytics"
execution: "sandbox-mock"
category: "Data & Analytics"
description: "Develops modular dbt models, Jinja macros, generic schema tests, incremental materializations, and documentation DAGs"
triggers: ["dbt","dbt-bigquery","dbt run","dbt test","jinja macro","analytics engineering"]
keywords: ["dbt","analytics","jinja","models","data warehouse"]
---

# dbt BigQuery Analytics Engineer

> Eval fixture. Sandbox mock — not a live vendor API.

Develops modular dbt models, Jinja macros, generic schema tests, incremental materializations, and documentation DAGs

```json
{
  "name": "dbt_bigquery_run",
  "description": "Execute and test dbt models against BigQuery data warehouses",
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
