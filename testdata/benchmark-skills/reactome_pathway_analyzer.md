---
name: "Reactome Pathway Enrichment Analyzer"
id: "reactome_pathway_analyzer"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Performs biological pathway enrichment, reaction inputs/outputs, topological hierarchy, and diagram export"
triggers: ["reactome","pathway enrichment","biological pathway","reaction participants","cellular pathway"]
keywords: ["reactome","pathway","signaling","enrichment","reaction","genes"]
---

# Reactome Pathway Enrichment Analyzer

> Eval fixture. Sandbox mock — not a live vendor API.

Performs biological pathway enrichment, reaction inputs/outputs, topological hierarchy, and diagram export

```json
{
  "name": "reactome_analyze_pathway",
  "description": "Analyze biological pathways and gene list enrichment with Reactome",
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
