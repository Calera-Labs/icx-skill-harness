---
name: "ChEMBL Drug Target & Bioactivity Database"
id: "chembl_target_bioactivity"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Queries ChEMBL for bioactive small molecules, IC50/Ki target affinities, approved drug mechanisms, and clinical phase"
triggers: ["chembl","ic50","ki affinity","drug target","bioactivity","target inhibition"]
keywords: ["chembl","pharma","ic50","target","drug","affinity"]
---

# ChEMBL Drug Target & Bioactivity Database

> Eval fixture. Sandbox mock — not a live vendor API.

Queries ChEMBL for bioactive small molecules, IC50/Ki target affinities, approved drug mechanisms, and clinical phase

```json
{
  "name": "chembl_target_query",
  "description": "Query ChEMBL for drug targets, IC50/Ki bioactivities, and approved drugs",
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
