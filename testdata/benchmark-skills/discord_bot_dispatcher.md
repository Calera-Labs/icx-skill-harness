---
name: "Discord Webhook & Bot Dispatcher"
id: "discord_bot_dispatcher"
execution: "sandbox-mock"
category: "Communication"
description: "Dispatches rich Discord embed cards, monitors server channels, and executes bot interactions"
triggers: ["discord","discord webhook","discord embed","discord bot"]
keywords: ["discord","bot","embed","webhook","community"]
---

# Discord Webhook & Bot Dispatcher

> Eval fixture. Sandbox mock — not a live vendor API.

Dispatches rich Discord embed cards, monitors server channels, and executes bot interactions

```json
{
  "name": "discord_dispatch_embed",
  "description": "Send rich embed messages to Discord channels",
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
