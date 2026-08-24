---
name: "Plaid Open Banking & ACH Verifier"
id: "plaid_bank_account_link"
execution: "sandbox-mock"
category: "FinTech"
description: "Verifies bank account routing, real-time balances, ACH identity tokens, and transaction history via Plaid"
triggers: ["plaid","bank verification","ach routing","bank balance","open banking"]
keywords: ["plaid","banking","ach","balance","transactions","fintech"]
---

# Plaid Open Banking & ACH Verifier

> Eval fixture. Sandbox mock — not a live vendor API.

Verifies bank account routing, real-time balances, ACH identity tokens, and transaction history via Plaid

```json
{
  "name": "plaid_verify_account",
  "description": "Verify bank accounts and check real-time balances via Plaid",
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
