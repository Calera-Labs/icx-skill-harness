---
name: "Elastic APM Performance Monitor"
id: "elastic_apm_collector"
execution: "sandbox-mock"
category: "Observability"
description: "Tracks transaction spans and unhandled exceptions in Elastic APM"
triggers: ["elastic apm","spans","transactions","elastic trace"]
keywords: ["elastic","apm","spans","transactions","monitoring"]
---

# Elastic APM Performance Monitor

> Eval fixture. Sandbox mock — not a live vendor API.

Tracks transaction spans and unhandled exceptions in Elastic APM

```json
{
  "name": "elastic_apm_query",
  "description": "Query Elastic APM spans and errors",
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
