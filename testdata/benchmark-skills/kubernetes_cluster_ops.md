---
name: "Kubernetes Cluster Orchestrator"
id: "kubernetes_cluster_ops"
execution: "sandbox-mock"
category: "DevOps"
description: "Scales pods, updates deployments, and manages namespaces in Kubernetes clusters"
triggers: ["kubernetes","k8s","kubectl","replicas","scale","namespace","deployment"]
keywords: ["k8s","kubernetes","pod","deployment","kubectl","cluster"]
---

# Kubernetes Cluster Orchestrator

> Eval fixture. Sandbox mock — not a live vendor API.

Scales pods, updates deployments, and manages namespaces in Kubernetes clusters

```json
{
  "name": "kubectl_orchestrator",
  "description": "Execute kubectl operations on Kubernetes clusters",
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
