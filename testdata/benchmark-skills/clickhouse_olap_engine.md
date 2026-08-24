---
name: "ClickHouse Real-Time OLAP Engine"
id: "clickhouse_olap_engine"
execution: "sandbox-mock"
category: "Database"
description: "Runs real-time vector and analytical queries on billions of events"
triggers: ["clickhouse","olap","columnar","realtime","event analytics"]
keywords: ["clickhouse","olap","columnar","analytics","sql"]
---

# ClickHouse Real-Time OLAP Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Runs real-time vector and analytical queries on billions of events

```json
{
  "name": "clickhouse_sql_query",
  "description": "Query ClickHouse columnar database",
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
