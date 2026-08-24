---
name: "gnomAD Population Allele Frequency"
id: "gnomad_allele_frequency"
execution: "sandbox-mock"
category: "Genomics"
description: "Queries the Genome Aggregation Database for population allele frequencies and gene constraint metrics (pLI, LOEUF)"
triggers: ["gnomad","allele frequency","pli constraint","loeuf","population genomics","loss of function"]
keywords: ["gnomad","population","allele","frequency","constraint","genomics"]
---

# gnomAD Population Allele Frequency

> Eval fixture. Sandbox mock — not a live vendor API.

Queries the Genome Aggregation Database for population allele frequencies and gene constraint metrics (pLI, LOEUF)

```json
{
  "name": "gnomad_query_frequency",
  "description": "Query gnomAD for population variant frequencies and gene constraint metrics",
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
