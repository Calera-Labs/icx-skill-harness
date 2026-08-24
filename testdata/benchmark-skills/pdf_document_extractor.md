---
name: "PDF Text & Layout Extractor"
id: "pdf_document_extractor"
execution: "sandbox-mock"
category: "Document Processing"
description: "Extracts text, structured tables, metadata, and form fields from PDF documents with layout preservation"
triggers: ["pdf","ocr","extract pdf","acrobat","pdf table","pdf layout"]
keywords: ["pdf","extraction","tables","ocr","document","reader"]
---

# PDF Text & Layout Extractor

> Eval fixture. Sandbox mock — not a live vendor API.

Extracts text, structured tables, metadata, and form fields from PDF documents with layout preservation

```json
{
  "name": "pdf_extractor",
  "description": "Extract text, tables, and structured data from PDF files",
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
