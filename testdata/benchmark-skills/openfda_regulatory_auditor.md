---
name: "openFDA Regulatory & Safety Auditor"
id: "openfda_regulatory_auditor"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Queries openFDA API for adverse drug event reports, 510(k) medical device clearances, NDC numbers, and recalls"
triggers: ["openfda","fda adverse event","510k","ndc lookup","drug recall","fda approval"]
keywords: ["fda","openfda","safety","adverse","pharma","clearance"]
---

# openFDA Regulatory & Safety Auditor

> Eval fixture. Sandbox mock — not a live vendor API.

Queries openFDA API for adverse drug event reports, 510(k) medical device clearances, NDC numbers, and recalls

```json
{
  "name": "openfda_query_records",
  "description": "Query openFDA database for drug safety, recalls, and device approvals",
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
