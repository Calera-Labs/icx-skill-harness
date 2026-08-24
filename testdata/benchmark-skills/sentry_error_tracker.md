---
name: "Sentry Real-Time Error Tracker"
id: "sentry_error_tracker"
execution: "sandbox-mock"
category: "Observability"
description: "Retrieves stack traces, unhandled exceptions, and breadcrumbs from Sentry"
triggers: ["sentry","stacktrace","error","exception","sentry issue"]
keywords: ["sentry","error","stacktrace","exception","breadcrumbs"]
---

# Sentry Real-Time Error Tracker

> Eval fixture. Sandbox mock — not a live vendor API.

Retrieves stack traces, unhandled exceptions, and breadcrumbs from Sentry

```json
{
  "name": "sentry_get_issue",
  "description": "Retrieve issue details and stack traces from Sentry",
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
