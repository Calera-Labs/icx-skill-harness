---
name: "InterPro Protein Family & Pfam Architect"
id: "interpro_domain_architect"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Scans protein sequences for functional domains, Pfam families, active sites, and deep-learning InterPro-N models"
triggers: ["interpro","pfam","protein domain","domain architecture","interpro-n","protein family"]
keywords: ["interpro","pfam","domain","family","protein","signature"]
---

# InterPro Protein Family & Pfam Architect

> Eval fixture. Sandbox mock — not a live vendor API.

Scans protein sequences for functional domains, Pfam families, active sites, and deep-learning InterPro-N models

```json
{
  "name": "interpro_domain_scan",
  "description": "Scan protein sequences against InterPro and Pfam domain signatures",
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
