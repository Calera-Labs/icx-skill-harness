---
name: "PowerPoint Deck Builder & Stylist"
id: "pptx_presentation_creator"
execution: "sandbox-mock"
category: "Document Processing"
description: "Creates professional Microsoft PowerPoint .pptx slide presentations with layouts, bullets, and speaker notes"
triggers: ["pptx","powerpoint","slides","deck","presentation","slide layout"]
keywords: ["pptx","presentation","slideshow","bullet","keynote"]
---

# PowerPoint Deck Builder & Stylist

> Eval fixture. Sandbox mock — not a live vendor API.

Creates professional Microsoft PowerPoint .pptx slide presentations with layouts, bullets, and speaker notes

```json
{
  "name": "pptx_deck_builder",
  "description": "Generate and format PowerPoint .pptx presentation decks",
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
