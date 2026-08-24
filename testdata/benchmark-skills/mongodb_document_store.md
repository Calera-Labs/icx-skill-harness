---
name: "MongoDB NoSQL Document Store"
id: "mongodb_document_store"
execution: "sandbox-mock"
category: "Database"
description: "Runs aggregation pipelines and CRUD operations on MongoDB collections"
triggers: ["mongodb","nosql","document","mongo","aggregation pipeline"]
keywords: ["mongodb","mongo","nosql","bson","document"]
---

# MongoDB NoSQL Document Store

> Eval fixture. Sandbox mock — not a live vendor API.

Runs aggregation pipelines and CRUD operations on MongoDB collections

```json
{
  "name": "mongodb_aggregate",
  "description": "Execute MongoDB aggregation pipelines",
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
