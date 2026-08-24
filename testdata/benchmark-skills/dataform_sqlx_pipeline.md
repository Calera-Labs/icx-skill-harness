---
name: "Dataform SQLX Pipeline Engineer"
id: "dataform_sqlx_pipeline"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Authors Dataform SQLX pipeline definitions, assertions, incremental table materializations, and workflow_settings.yaml"
triggers: ["dataform","sqlx","dataform pipeline","assertions","workflow_settings"]
keywords: ["dataform","sqlx","elt","bigquery","transformation"]
---

# Dataform SQLX Pipeline Engineer

> Eval fixture. Sandbox mock — not a live vendor API.

Authors Dataform SQLX pipeline definitions, assertions, incremental table materializations, and workflow_settings.yaml

```json
{
  "name": "dataform_compile_run",
  "description": "Compile and execute Dataform SQLX transformation pipelines",
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
