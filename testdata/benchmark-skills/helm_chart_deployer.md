---
name: "Helm Kubernetes Chart Deployer"
id: "helm_chart_deployer"
execution: "sandbox-mock"
category: "DevOps"
description: "Renders Helm templates and installs chart releases into Kubernetes"
triggers: ["helm","chart","kubernetes deploy","release","helm install"]
keywords: ["helm","chart","values","release","k8s"]
---

# Helm Kubernetes Chart Deployer

> Eval fixture. Sandbox mock — not a live vendor API.

Renders Helm templates and installs chart releases into Kubernetes

```json
{
  "name": "helm_install_release",
  "description": "Install or upgrade Helm chart releases",
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
