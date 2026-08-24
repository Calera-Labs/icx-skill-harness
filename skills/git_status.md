---
name: "Git Status"
id: "git_status"
category: "Starter"
execution: "sandbox-mock"
description: "Sandbox mock of a git status skill. This harness does not run git."
triggers: ["git status", "working tree", "repo status"]
keywords: ["git", "status", "diff", "mock"]
---

# Git Status

Default tool execution in this repo is a sandbox mock. It does not inspect a real repository.

```json
{
  "name": "git_status",
  "description": "Sandbox mock of git working tree status",
  "parameters": {
    "type": "object",
    "properties": {
      "repo_path": {"type": "string", "description": "Path to a git repository (not read)"}
    },
    "required": ["repo_path"]
  }
}
```
