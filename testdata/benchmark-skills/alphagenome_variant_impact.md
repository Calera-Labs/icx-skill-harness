---
name: "AlphaGenome Non-Coding Regulatory Predictor"
id: "alphagenome_variant_impact"
execution: "sandbox-mock"
category: "Genomics"
description: "Predicts non-coding variant impact on gene expression, chromatin accessibility (DNase), and histone marks"
triggers: ["alphagenome","non-coding variant","chromatin accessibility","dnase","regulatory variant effect"]
keywords: ["alphagenome","genomics","expression","noncoding","chromatin"]
---

# AlphaGenome Non-Coding Regulatory Predictor

> Eval fixture. Sandbox mock — not a live vendor API.

Predicts non-coding variant impact on gene expression, chromatin accessibility (DNase), and histone marks

```json
{
  "name": "alphagenome_variant_predict",
  "description": "Predict regulatory and expression effects of non-coding variants with AlphaGenome",
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
