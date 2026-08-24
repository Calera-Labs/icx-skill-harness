---
name: "QuickGO Gene Ontology Term Navigator"
id: "quickgo_ontology_mapper"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Maps Gene Ontology (GO) terms across Biological Process, Molecular Function, and Cellular Component hierarchies"
triggers: ["quickgo","gene ontology","go term","biological process","molecular function","eco code"]
keywords: ["quickgo","go","ontology","biological process","molecular function"]
---

# QuickGO Gene Ontology Term Navigator

> Eval fixture. Sandbox mock — not a live vendor API.

Maps Gene Ontology (GO) terms across Biological Process, Molecular Function, and Cellular Component hierarchies

```json
{
  "name": "quickgo_fetch_terms",
  "description": "Query QuickGO for Gene Ontology terms and ECO evidence codes",
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
