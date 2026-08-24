---
name: "Word Document Architect & Formatter"
id: "docx_document_architect"
execution: "sandbox-mock"
category: "Document Processing"
description: "Inspects, generates, formats and modifies Microsoft Word .docx documents, paragraph styles, headings and tables"
triggers: ["docx","word","document","paragraph","heading","formatting","word document"]
keywords: ["docx","document","office","style","table","msword"]
---

# Word Document Architect & Formatter

> Eval fixture. Sandbox mock — not a live vendor API.

Inspects, generates, formats and modifies Microsoft Word .docx documents, paragraph styles, headings and tables

```json
{
  "name": "docx_manipulator",
  "description": "Create and format Microsoft Word .docx documents",
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
