---
name: "Twilio SMS & 2FA Auth Dispatcher"
id: "twillio_sms_notifier"
execution: "sandbox-mock"
category: "Communication"
description: "Dispatches SMS verification codes and 2FA alerts via Twilio API"
triggers: ["twilio","sms","2fa","text message","send sms"]
keywords: ["twilio","sms","2fa","text","phone"]
---

# Twilio SMS & 2FA Auth Dispatcher

> Eval fixture. Sandbox mock — not a live vendor API.

Dispatches SMS verification codes and 2FA alerts via Twilio API

```json
{
  "name": "twilio_send_sms",
  "description": "Send SMS text messages via Twilio API",
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
