---
name: "IAM Security & Token Rotator"
id: "iam_security_guard"
execution: "sandbox-mock"
category: "Security"
description: "Rotates OAuth2 service account tokens, audits access logs, and validates HMAC signatures"
triggers: ["iam","oauth2","token","rotate","security","credentials"]
keywords: ["iam","oauth2","token","security","auth"]
---

# IAM Security & Token Rotator

> Eval fixture. Sandbox mock — not a live vendor API.

Rotates OAuth2 service account tokens, audits access logs, and validates HMAC signatures

```json
{
  "name": "iam_auth_rotator",
  "description": "Rotate IAM credentials and OAuth tokens",
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
