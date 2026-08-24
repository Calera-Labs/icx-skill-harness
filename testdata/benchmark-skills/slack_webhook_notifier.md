---
name: "Slack Webhook & Alert Bot"
id: "slack_webhook_notifier"
execution: "sandbox-mock"
category: "Communication"
description: "Dispatches rich Block Kit notifications and incident alerts to Slack"
triggers: ["slack","notification","alert","channel","message","slack webhook"]
keywords: ["slack","message","webhook","block kit","channel"]
---

# Slack Webhook & Alert Bot

> Eval fixture. Sandbox mock — not a live vendor API.

Dispatches rich Block Kit notifications and incident alerts to Slack

```json
{
  "name": "slack_send_message",
  "description": "Send rich notification messages to Slack channels",
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
