---
name: "AlphaFold Structure & Confidence Analyzer"
id: "alphafold_structure_predictor"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Retrieves AlphaFold 3D protein structure predictions, per-residue pLDDT confidence scores, and domain boundaries"
triggers: ["alphafold","uniprot structure","plddt","protein structure","domain boundary","alphafold database"]
keywords: ["alphafold","protein","plddt","structure","uniprot","biology"]
---

# AlphaFold Structure & Confidence Analyzer

> Eval fixture. Sandbox mock — not a live vendor API.

Retrieves AlphaFold 3D protein structure predictions, per-residue pLDDT confidence scores, and domain boundaries

```json
{
  "name": "alphafold_fetch_analyze",
  "description": "Fetch and analyze AlphaFold 3D protein structures and pLDDT scores",
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
