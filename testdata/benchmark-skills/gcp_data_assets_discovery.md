---
name: "Dataplex & GCP Data Asset Discoverer"
id: "gcp_data_assets_discovery"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Searches enterprise data catalogs across BigQuery, BigLake, Spanner, and Dataplex with schema and metadata profiling"
triggers: ["dataplex","data assets","discover tables","data catalog","spanner table","table schema","bigquery datasets","discover data assets","gcp data assets"]
keywords: ["dataplex","catalog","assets","schema","discovery","bigquery","datasets"]
---

# Dataplex & GCP Data Asset Discoverer

> Eval fixture. Sandbox mock — not a live vendor API.

Searches enterprise data catalogs across BigQuery, BigLake, Spanner, and Dataplex with schema and metadata profiling

```json
{
  "name": "dataplex_asset_discover",
  "description": "Discover data assets, tables, and schemas in Dataplex and BigQuery",
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
