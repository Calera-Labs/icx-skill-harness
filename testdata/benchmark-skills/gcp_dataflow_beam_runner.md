---
name: "Apache Beam & Dataflow Runner"
id: "gcp_dataflow_beam_runner"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Authors Apache Beam Java/Python pipelines, packages Flex Templates, and troubleshoots streaming autoscaling bottlenecks"
triggers: ["dataflow","apache beam","beam pipeline","flex template","streaming job"]
keywords: ["dataflow","beam","streaming","batch","gcp"]
---

# Apache Beam & Dataflow Runner

> Eval fixture. Sandbox mock — not a live vendor API.

Authors Apache Beam Java/Python pipelines, packages Flex Templates, and troubleshoots streaming autoscaling bottlenecks

```json
{
  "name": "dataflow_job_manager",
  "description": "Deploy and inspect Apache Beam pipelines on Google Cloud Dataflow",
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
