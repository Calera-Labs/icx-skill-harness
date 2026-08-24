---
name: "Linear Issue & Project Tracker"
id: "linear_project_tracker"
execution: "sandbox-mock"
category: "Management"
description: "Creates Linear issue tickets, assigns cycle sprints, updates priority states, and manages project milestones"
triggers: ["linear","linear issue","sprint cycle","project roadmap","ticket status"]
keywords: ["linear","issue","sprint","project","cycle","tracker"]
---

# Linear Issue & Project Tracker

> Eval fixture. Sandbox mock — not a live vendor API.

Creates Linear issue tickets, assigns cycle sprints, updates priority states, and manages project milestones

```json
{
  "name": "linear_manage_issue",
  "description": "Create and update Linear project issues and cycle sprints",
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
