---
name: "Argo Workflows Cloud-Native Engine"
id: "argo_workflows_cicd"
execution: "sandbox-mock"
category: "DevOps"
description: "Orchestrates DAG-based container workflows on Kubernetes clusters"
triggers: ["argo","workflows","dag","kubernetes workflow","argo submit"]
keywords: ["argo","workflows","dag","kubernetes","ci"]
---

# Argo Workflows Cloud-Native Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Orchestrates DAG-based container workflows on Kubernetes clusters

```json
{
  "name": "argo_submit_workflow",
  "description": "Submit and monitor Argo Workflow DAGs",
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
