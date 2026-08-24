---
name: "SQL Query Optimizer & Plan Tuner"
id: "sql_optimizer_query_tuner"
execution: "sandbox-mock"
category: "Database"
description: "Analyzes EXPLAIN ANALYZE execution plans, suggests index tuning, cost-based join reordering, and subquery flattening"
triggers: ["sql optimizer","explain analyze","query plan","index tuning","slow query","query performance"]
keywords: ["sql","query","optimizer","explain","index","execution plan"]
---

# SQL Query Optimizer & Plan Tuner

> Eval fixture. Sandbox mock — not a live vendor API.

Analyzes EXPLAIN ANALYZE execution plans, suggests index tuning, cost-based join reordering, and subquery flattening

```json
{
  "name": "sql_plan_optimize",
  "description": "Optimize SQL execution plans and index strategies",
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
