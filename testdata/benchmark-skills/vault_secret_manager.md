---
name: "HashiCorp Vault Secret Manager"
id: "vault_secret_manager"
execution: "sandbox-mock"
category: "Security"
description: "Reads and writes dynamic secrets and encryption keys from Vault"
triggers: ["vault","secret","hashicorp","encryption","vault secret"]
keywords: ["vault","secret","hashicorp","token","encryption"]
---

# HashiCorp Vault Secret Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Reads and writes dynamic secrets and encryption keys from Vault

```json
{
  "name": "vault_secret_read",
  "description": "Read and write HashiCorp Vault secrets",
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
