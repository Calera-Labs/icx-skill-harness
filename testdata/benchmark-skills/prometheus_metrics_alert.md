---
name: "Prometheus Metrics & Alert Manager"
id: "prometheus_metrics_alert"
execution: "sandbox-mock"
category: "Observability"
description: "Queries Prometheus metrics, calculates P99 latency, and triggers alerts"
triggers: ["prometheus","metrics","p99","latency","alert","promql"]
keywords: ["prometheus","promql","metrics","p99","alerts"]
---

# Prometheus Metrics & Alert Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Queries Prometheus metrics, calculates P99 latency, and triggers alerts

```json
{
  "name": "prometheus_query",
  "description": "Query Prometheus time-series metrics",
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
