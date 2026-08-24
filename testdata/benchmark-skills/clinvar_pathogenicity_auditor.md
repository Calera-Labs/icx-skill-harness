---
name: "ClinVar Genomic Pathogenicity Auditor"
id: "clinvar_pathogenicity_auditor"
execution: "sandbox-mock"
category: "Genomics"
description: "Evaluates human genetic variants against ClinVar for Pathogenic, Likely Pathogenic, Benign, or VUS classifications"
triggers: ["clinvar","pathogenicity","variant clinical significance","vus","pathogenic variant","acmg"]
keywords: ["clinvar","genomics","pathogenic","variant","mutation","clinical"]
---

# ClinVar Genomic Pathogenicity Auditor

> Eval fixture. Sandbox mock — not a live vendor API.

Evaluates human genetic variants against ClinVar for Pathogenic, Likely Pathogenic, Benign, or VUS classifications

```json
{
  "name": "clinvar_variant_lookup",
  "description": "Look up human genomic variant pathogenicity and clinical evidence in ClinVar",
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
