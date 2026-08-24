---
name: "GitHub Actions CI/CD Pipeline Manager"
id: "github_actions_ci"
execution: "sandbox-mock"
category: "DevOps"
description: "Triggers workflow dispatches and monitors CI test runners"
triggers: ["github actions","ci","cd","pipeline","workflow dispatch"]
keywords: ["github","actions","ci","cd","workflow","runner"]
---

# GitHub Actions CI/CD Pipeline Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Triggers workflow dispatches and monitors CI test runners

```json
{
  "name": "github_actions_trigger",
  "description": "Trigger and monitor GitHub Actions workflows",
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
