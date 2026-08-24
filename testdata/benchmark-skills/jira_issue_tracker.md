---
name: "Jira Issue & Sprint Board Manager"
id: "jira_issue_tracker"
execution: "sandbox-mock"
category: "Management"
description: "Creates tickets, updates sprint story points, and assigns bugs in Jira"
triggers: ["jira","ticket","sprint","bug","issue","jira ticket"]
keywords: ["jira","ticket","sprint","bug","issue"]
---

# Jira Issue & Sprint Board Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Creates tickets, updates sprint story points, and assigns bugs in Jira

```json
{
  "name": "jira_update_ticket",
  "description": "Update Jira tickets and sprint boards",
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
