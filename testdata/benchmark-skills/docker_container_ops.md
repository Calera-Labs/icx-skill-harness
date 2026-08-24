---
name: "Docker Container Lifecycle Manager"
id: "docker_container_ops"
execution: "sandbox-mock"
category: "DevOps"
description: "Manages docker containers, inspects container logs, and restarts microservices"
triggers: ["docker","container","restart","container health","microservice"]
keywords: ["docker","container","image","daemon","dockerfile"]
---

# Docker Container Lifecycle Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Manages docker containers, inspects container logs, and restarts microservices

```json
{
  "name": "docker_manager",
  "description": "Inspect and manage Docker container lifecycle",
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
