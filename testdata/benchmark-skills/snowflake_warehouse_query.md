---
name: "Snowflake Cloud Data Warehouse"
id: "snowflake_warehouse_query"
execution: "sandbox-mock"
category: "Data"
description: "Executes analytics SQL on Snowflake warehouses with auto-scaling"
triggers: ["snowflake","warehouse","analytics","data warehouse","snowpark"]
keywords: ["snowflake","warehouse","sql","data","cloud"]
---

# Snowflake Cloud Data Warehouse

> Eval fixture. Sandbox mock — not a live vendor API.

Executes analytics SQL on Snowflake warehouses with auto-scaling

```json
{
  "name": "snowflake_query_exec",
  "description": "Execute analytical queries on Snowflake",
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
