---
name: "Sandbox Echo"
id: "echo_sandbox"
category: "Starter"
execution: "sandbox-mock"
description: "Echoes arguments."
triggers: ["echo", "sandbox", "ping skill", "hello harness"]
keywords: ["echo", "sandbox", "starter"]
---

# Sandbox Echo

```json
{
  "name": "echo_sandbox",
  "description": "Echo arguments.",
  "parameters": {
    "type": "object",
    "properties": {
      "message": {"type": "string", "description": "Text to echo"}
    },
    "required": ["message"]
  }
}
```
