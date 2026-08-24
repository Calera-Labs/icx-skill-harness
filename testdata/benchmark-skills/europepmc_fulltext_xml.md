---
name: "Europe PMC Full-Text & Bio-Entities"
id: "europepmc_fulltext_xml"
execution: "sandbox-mock"
category: "Research & Literature"
description: "Retrieves open-access biomedical literature in full-text XML and identifies mined bio-entities and chemicals"
triggers: ["europe pmc","europepmc","full text xml","pmcid","biomedical fulltext"]
keywords: ["europepmc","pmc","xml","fulltext","biomedical"]
---

# Europe PMC Full-Text & Bio-Entities

> Eval fixture. Sandbox mock — not a live vendor API.

Retrieves open-access biomedical literature in full-text XML and identifies mined bio-entities and chemicals

```json
{
  "name": "europepmc_fetch_fulltext",
  "description": "Fetch full-text XML articles and entity annotations from Europe PMC",
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
