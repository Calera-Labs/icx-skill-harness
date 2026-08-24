---
name: "DCF & Valuation Financial Modeler"
id: "val_dcf_modeler"
execution: "sandbox-mock"
category: "Finance"
description: "Calculates discounted cash flows, terminal value, and WACC hurdle rates with sensitivity tables"
triggers: ["dcf","valuation","wacc","cash flow","terminal value"]
keywords: ["dcf","valuation","wacc","discounted","finance","cash flow"]
---

# DCF & Valuation Financial Modeler

> Eval fixture. Sandbox mock — not a live vendor API.

Calculates discounted cash flows, terminal value, and WACC hurdle rates with sensitivity tables

```json
{
  "name": "valuation_dcf_calc",
  "description": "Compute DCF financial models and valuation",
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
