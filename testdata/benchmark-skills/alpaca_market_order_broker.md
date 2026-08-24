---
name: "Alpaca Market Order Broker"
id: "alpaca_market_order_broker"
execution: "sandbox-mock"
category: "FinTech"
description: "Executes equity and options orders, inspects portfolio buying power, and sets bracket stop-loss orders on Alpaca"
triggers: ["alpaca","market order","limit order","bracket order","buying power","stock trade"]
keywords: ["alpaca","trading","stocks","orders","portfolio","broker"]
---

# Alpaca Market Order Broker

> Eval fixture. Sandbox mock — not a live vendor API.

Executes equity and options orders, inspects portfolio buying power, and sets bracket stop-loss orders on Alpaca

```json
{
  "name": "alpaca_place_order",
  "description": "Submit and manage equity market and limit orders via Alpaca",
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
