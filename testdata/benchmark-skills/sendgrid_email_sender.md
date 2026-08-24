---
name: "SendGrid Transactional Email Sender"
id: "sendgrid_email_sender"
execution: "sandbox-mock"
category: "Communication"
description: "Dispatches password resets and transactional email receipts via SendGrid"
triggers: ["sendgrid","email","transactional email","receipt email","send email"]
keywords: ["sendgrid","email","smtp","transactional","mail"]
---

# SendGrid Transactional Email Sender

> Eval fixture. Sandbox mock — not a live vendor API.

Dispatches password resets and transactional email receipts via SendGrid

```json
{
  "name": "sendgrid_send_email",
  "description": "Send transactional emails via SendGrid API",
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
