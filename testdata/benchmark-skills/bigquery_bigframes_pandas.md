---
name: "BigQuery DataFrames & BigFrames"
id: "bigquery_bigframes_pandas"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Executes pandas and scikit-learn DataFrame operations on petabyte datasets backed by BigQuery compute engine"
triggers: ["bigframes","bigquery dataframe","pandas bigquery","scikit-learn bigquery"]
keywords: ["bigframes","pandas","dataframe","python","bigquery"]
---

# BigQuery DataFrames & BigFrames

> Eval fixture. Sandbox mock — not a live vendor API.

Executes pandas and scikit-learn DataFrame operations on petabyte datasets backed by BigQuery compute engine

```json
{
  "name": "bigframes_dataframe_ops",
  "description": "Execute BigFrames pandas operations on BigQuery tables",
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
