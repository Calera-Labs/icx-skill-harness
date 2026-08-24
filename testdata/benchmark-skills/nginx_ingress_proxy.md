---
name: "NGINX Ingress & Reverse Proxy"
id: "nginx_ingress_proxy"
execution: "sandbox-mock"
category: "Networking"
description: "Configures upstream routing, SSL termination, and rate limiting rules"
triggers: ["nginx","proxy","ssl","reverse proxy","ingress"]
keywords: ["nginx","proxy","reverse","ssl","upstream"]
---

# NGINX Ingress & Reverse Proxy

> Eval fixture. Sandbox mock — not a live vendor API.

Configures upstream routing, SSL termination, and rate limiting rules

```json
{
  "name": "nginx_reload_config",
  "description": "Configure and reload NGINX reverse proxy",
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
