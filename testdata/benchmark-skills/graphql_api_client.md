---
name: "GraphQL Schema & Query Executor"
id: "graphql_api_client"
execution: "sandbox-mock"
category: "API"
description: "Introspects GraphQL schemas and executes GraphQL mutations and queries"
triggers: ["graphql","mutation","graphql query","graphql schema"]
keywords: ["graphql","query","mutation","schema","api"]
---

# GraphQL Schema & Query Executor

> Eval fixture. Sandbox mock — not a live vendor API.

Introspects GraphQL schemas and executes GraphQL mutations and queries

```json
{
  "name": "graphql_execute",
  "description": "Execute GraphQL queries and mutations",
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
