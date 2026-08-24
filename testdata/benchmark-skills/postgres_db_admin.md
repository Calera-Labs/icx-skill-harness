---
name: "PostgreSQL Database Administrator"
id: "postgres_db_admin"
execution: "sandbox-mock"
category: "Database"
description: "Executes parameterized SQL queries, schema migrations, and transaction management"
triggers: ["postgres","sql","database","table","transactions table","settlement status","committed","sql query"]
keywords: ["postgres","postgresql","sql","database","query"]
---

# PostgreSQL Database Administrator

> Eval fixture. Sandbox mock — not a live vendor API.

Executes parameterized SQL queries, schema migrations, and transaction management

```json
{
  "name": "postgres_executor",
  "description": "Execute SQL queries against PostgreSQL database",
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
