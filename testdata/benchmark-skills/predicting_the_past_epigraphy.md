---
name: "Ancient Epigraphy & Text Restorer"
id: "predicting_the_past_epigraphy"
execution: "sandbox-mock"
category: "Humanities & AI"
description: "Restores missing characters, dates, and attributes ancient Greek and Latin inscriptions using Aeneas and Ithaca models"
triggers: ["epigraphy","ancient greek","latin inscription","ithaca","aeneas","text restoration"]
keywords: ["epigraphy","ancient","greek","latin","history","restoration"]
---

# Ancient Epigraphy & Text Restorer

> Eval fixture. Sandbox mock — not a live vendor API.

Restores missing characters, dates, and attributes ancient Greek and Latin inscriptions using Aeneas and Ithaca models

```json
{
  "name": "epigraphy_restore_text",
  "description": "Restore, date, and locate ancient epigraphic texts using Aeneas/Ithaca",
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
