---
name: "SQLite Embedded Ledger Auditor"
id: "sqlite_local_auditor"
execution: "sandbox-mock"
category: "Database"
description: "Audits local ACID SQLite databases and verifies table indexes"
triggers: ["sqlite","embedded","ledger","integrity","sqlite db"]
keywords: ["sqlite","embedded","acid","database","local"]
---

# SQLite Embedded Ledger Auditor

> Eval fixture. Sandbox mock — not a live vendor API.

Audits local ACID SQLite databases and verifies table indexes

```json
{
  "name": "sqlite_verify_integrity",
  "description": "Verify SQLite database integrity",
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
