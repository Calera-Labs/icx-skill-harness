---
name: "Git & AST Code Patcher"
id: "git_code_patcher"
execution: "sandbox-mock"
category: "DevOps"
description: "Inspects repository AST symbol trees and generates unified diff git patches"
triggers: ["git","diff","patch","function","ast","refactor"]
keywords: ["git","patch","diff","refactor","ast","code"]
---

# Git & AST Code Patcher

> Eval fixture. Sandbox mock — not a live vendor API.

Inspects repository AST symbol trees and generates unified diff git patches

```json
{
  "name": "git_diff_patcher",
  "description": "Synthesize git diff patch for code refactoring",
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
