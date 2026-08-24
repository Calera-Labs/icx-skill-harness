---
name: "ENCODE SCREEN cis-Regulatory Elements"
id: "encode_screen_ccre_query"
execution: "sandbox-mock"
category: "Genomics"
description: "Searches ENCODE Registry of candidate cis-Regulatory Elements (cCREs) across human cell lines and ChIP-seq peaks"
triggers: ["encode","ccre","screen","regulatory element","promoter enhancer","chip-seq peak"]
keywords: ["encode","ccre","screen","enhancer","promoter","epigenomics"]
---

# ENCODE SCREEN cis-Regulatory Elements

> Eval fixture. Sandbox mock — not a live vendor API.

Searches ENCODE Registry of candidate cis-Regulatory Elements (cCREs) across human cell lines and ChIP-seq peaks

```json
{
  "name": "encode_ccre_search",
  "description": "Search candidate cis-regulatory elements in ENCODE SCREEN",
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
