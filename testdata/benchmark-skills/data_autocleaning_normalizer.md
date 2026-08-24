---
name: "Automated Data Quality & Normalizer"
id: "data_autocleaning_normalizer"
execution: "sandbox-mock"
category: "Data & Analytics"
description: "Applies automated data cleansing, schema mapping, null value imputation, duplicate removal, and type normalization"
triggers: ["data cleaning","schema mapping","data quality","null imputation","deduplication","data autocleaning"]
keywords: ["cleaning","quality","normalizer","schema","etl","imputation"]
---

# Automated Data Quality & Normalizer

> Eval fixture. Sandbox mock — not a live vendor API.

Applies automated data cleansing, schema mapping, null value imputation, duplicate removal, and type normalization

```json
{
  "name": "data_clean_transform",
  "description": "Execute automated data cleansing and schema normalization",
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
