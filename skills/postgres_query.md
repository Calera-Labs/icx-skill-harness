---
name: "Postgres Query"
id: "postgres_query"
category: "Starter"
execution: "sandbox-mock"
description: "Sandbox mock of a read-only SQL skill. This harness does not connect to Postgres."
triggers: ["postgres", "sql select", "database query"]
keywords: ["postgres", "sql", "dsn", "mock"]
---

# Postgres Query

Default tool execution in this repo is a sandbox mock. It does not read DATABASE_URL or run SQL.

```json
{
  "name": "postgres_query",
  "description": "Sandbox mock of a read-only SQL query skill",
  "parameters": {
    "type": "object",
    "properties": {
      "sql": {"type": "string", "description": "SELECT statement (not executed)"}
    },
    "required": ["sql"]
  }
}
```
