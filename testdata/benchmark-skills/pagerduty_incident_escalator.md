---
name: "PagerDuty On-Call Incident Escalator"
id: "pagerduty_incident_escalator"
execution: "sandbox-mock"
category: "DevOps"
description: "Triggers on-call alerts, creates high-urgency incidents, and manages schedules"
triggers: ["pagerduty","oncall","incident","escalation","pager","trigger incident"]
keywords: ["pagerduty","incident","oncall","escalate","alert"]
---

# PagerDuty On-Call Incident Escalator

> Eval fixture. Sandbox mock — not a live vendor API.

Triggers on-call alerts, creates high-urgency incidents, and manages schedules

```json
{
  "name": "pagerduty_trigger_incident",
  "description": "Trigger PagerDuty on-call incident alerts",
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
