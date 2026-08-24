---
name: "PubChem Chemical Compound Database"
id: "pubchem_cheminformatics"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Resolves chemical names, SMILES strings, PubChem CIDs, molecular weights, 2D/3D structures, and bioassays"
triggers: ["pubchem","smiles","cid","chemical structure","cheminformatics","molecular formula"]
keywords: ["pubchem","chemistry","smiles","compound","molecule","bioassay"]
---

# PubChem Chemical Compound Database

> Eval fixture. Sandbox mock — not a live vendor API.

Resolves chemical names, SMILES strings, PubChem CIDs, molecular weights, 2D/3D structures, and bioassays

```json
{
  "name": "pubchem_compound_search",
  "description": "Search PubChem for chemical compounds, SMILES, and bioactivity data",
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
