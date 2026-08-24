---
name: "Stripe Billing & Payment Reconciler"
id: "stripe_billing_ops"
execution: "sandbox-mock"
category: "FinTech"
description: "Reconciles customer invoices, subscription billing, payment intents, and webhooks"
triggers: ["stripe","invoice","payment","webhook","billing"]
keywords: ["stripe","billing","invoice","payment intent","subscription"]
---

# Stripe Billing & Payment Reconciler

> Eval fixture. Sandbox mock — not a live vendor API.

Reconciles customer invoices, subscription billing, payment intents, and webhooks

```json
{
  "name": "stripe_reconciler",
  "description": "Reconcile Stripe invoices and webhooks",
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
