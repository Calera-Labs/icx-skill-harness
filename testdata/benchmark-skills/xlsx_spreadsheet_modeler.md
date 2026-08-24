---
name: "Excel Spreadsheet & Formula Modeler"
id: "xlsx_spreadsheet_modeler"
execution: "sandbox-mock"
category: "Document Processing"
description: "Builds, audits, and recalculates Excel .xlsx financial spreadsheets with SUMIFS, XLOOKUP, pivot tables, and styling"
triggers: ["xlsx","excel","spreadsheet","formula","xlookup","pivot table","workbook"]
keywords: ["xlsx","excel","sheet","sumifs","vlookup","cells","financial model"]
---

# Excel Spreadsheet & Formula Modeler

> Eval fixture. Sandbox mock — not a live vendor API.

Builds, audits, and recalculates Excel .xlsx financial spreadsheets with SUMIFS, XLOOKUP, pivot tables, and styling

```json
{
  "name": "xlsx_sheet_modeler",
  "description": "Generate and audit Excel .xlsx financial spreadsheets and formulas",
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
