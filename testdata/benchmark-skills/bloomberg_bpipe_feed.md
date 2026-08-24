---
name: "Bloomberg B-PIPE Market Feed"
id: "bloomberg_bpipe_feed"
execution: "sandbox-mock"
category: "Finance"
description: "Streams institutional Level 2 market data, tick depth, swap rates, and yield curve spreads from Bloomberg B-PIPE"
triggers: ["bloomberg","bpipe","market depth","level 2 quote","yield spread","orderbook"]
keywords: ["bloomberg","market","quote","orderbook","bpipe","ticks"]
---

# Bloomberg B-PIPE Market Feed

> Eval fixture. Sandbox mock — not a live vendor API.

Streams institutional Level 2 market data, tick depth, swap rates, and yield curve spreads from Bloomberg B-PIPE

```json
{
  "name": "bloomberg_market_quote",
  "description": "Fetch real-time market depth and pricing from Bloomberg B-PIPE",
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
