---
name: "STRING Protein Interaction Network"
id: "string_ppi_network_analyzer"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Queries STRING database for physical and functional protein-protein interaction networks and confidence scores"
triggers: ["string database","protein interaction","ppi network","protein partners","interactome"]
keywords: ["string","ppi","interactome","protein","network","confidence"]
---

# STRING Protein Interaction Network

> Eval fixture. Sandbox mock — not a live vendor API.

Queries STRING database for physical and functional protein-protein interaction networks and confidence scores

```json
{
  "name": "string_ppi_query",
  "description": "Retrieve protein-protein interaction networks and confidence scores from STRING",
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
