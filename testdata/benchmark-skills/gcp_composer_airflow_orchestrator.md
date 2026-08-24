---
name: "Cloud Composer & Airflow Orchestrator"
id: "gcp_composer_airflow_orchestrator"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Orchestrates Apache Airflow DAGs on Cloud Composer (MSAA Gen 2 & 3), manages task dependencies, and debugs failures"
triggers: ["composer","airflow","dag","cloud composer","msaa","task dependency"]
keywords: ["composer","airflow","orchestration","dag","gcp"]
---

# Cloud Composer & Airflow Orchestrator

> Eval fixture. Sandbox mock — not a live vendor API.

Orchestrates Apache Airflow DAGs on Cloud Composer (MSAA Gen 2 & 3), manages task dependencies, and debugs failures

```json
{
  "name": "composer_dag_trigger",
  "description": "Trigger and monitor Apache Airflow DAGs in Cloud Composer",
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
