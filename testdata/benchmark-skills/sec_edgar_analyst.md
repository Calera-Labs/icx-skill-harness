---
name: "SEC EDGAR Financial Analyst"
id: "sec_edgar_analyst"
execution: "sandbox-mock"
category: "Finance"
description: "Retrieves 10-K, 10-Q, 8-K filings and exact GAAP metrics from SEC EDGAR (eval mock)"
triggers: ["sec","10-k","10-q","gaap","operating margin","edgar"]
keywords: ["sec","edgar","10k","10q","gaap","revenue","filings"]
---

# SEC EDGAR Financial Analyst

> Eval fixture. Sandbox mock — not a live vendor API.

Retrieves 10-K, 10-Q, 8-K filings and exact GAAP metrics from SEC EDGAR (eval mock)

```json
{
  "name": "sec_edgar_query",
  "description": "Query SEC EDGAR database for verified financial metrics",
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
