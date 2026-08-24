---
name: "Datadog APM & Trace Analyzer"
id: "datadog_apm_tracer"
execution: "sandbox-mock"
category: "Observability"
description: "Inspects distributed flame graphs and pinpoints trace bottlenecks"
triggers: ["datadog","apm","trace","flamegraph","datadog trace"]
keywords: ["datadog","apm","trace","flamegraph","spans"]
---

# Datadog APM & Trace Analyzer

> Eval fixture. Sandbox mock — not a live vendor API.

Inspects distributed flame graphs and pinpoints trace bottlenecks

```json
{
  "name": "datadog_trace_inspect",
  "description": "Inspect Datadog APM distributed traces",
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
