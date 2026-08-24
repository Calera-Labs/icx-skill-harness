---
name: "GitHub Pull Request Reviewer"
id: "github_pr_auditor"
execution: "sandbox-mock"
category: "DevOps"
description: "Audits PR changes, runs static analysis, and submits review comments"
triggers: ["github","pr","pull request","review","pr diff"]
keywords: ["github","pr","review","diff","pull request"]
---

# GitHub Pull Request Reviewer

> Eval fixture. Sandbox mock — not a live vendor API.

Audits PR changes, runs static analysis, and submits review comments

```json
{
  "name": "github_pr_review",
  "description": "Review and audit GitHub pull requests",
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
