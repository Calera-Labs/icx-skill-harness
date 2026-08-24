---
name: "Cloudflare DNS & WAF Rule Manager"
id: "cloudflare_dns_waf"
execution: "sandbox-mock"
category: "Networking"
description: "Updates DNS A/CNAME records and configures Cloudflare WAF firewall rules"
triggers: ["cloudflare","dns","waf","firewall","dns record"]
keywords: ["cloudflare","dns","waf","firewall","cdn"]
---

# Cloudflare DNS & WAF Rule Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Updates DNS A/CNAME records and configures Cloudflare WAF firewall rules

```json
{
  "name": "cloudflare_update_dns",
  "description": "Update Cloudflare DNS records and WAF",
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
