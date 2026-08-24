---
name: "Grafana Dashboard & Alert Manager"
id: "grafana_dashboard_manager"
execution: "sandbox-mock"
category: "Observability"
description: "Provisions Grafana dashboard JSON models, configures panel alert thresholds, and links data sources"
triggers: ["grafana","dashboard json","grafana panel","alert threshold"]
keywords: ["grafana","dashboard","panel","visualization","alerts"]
---

# Grafana Dashboard & Alert Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Provisions Grafana dashboard JSON models, configures panel alert thresholds, and links data sources

```json
{
  "name": "grafana_manage_dashboard",
  "description": "Provision and update Grafana dashboards and alerts",
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
