---
name: "Code Review & Security Hardener"
id: "code_review_security_auditor"
execution: "sandbox-mock"
category: "DevOps & Security"
description: "Audits source code for OWASP Top 10 vulnerabilities, SQL injection, XSS, buffer overflows, and architectural smells"
triggers: ["code review","security audit","owasp","sql injection","vulnerability scan","security hardener"]
keywords: ["security","audit","cve","owasp","taint","code review"]
---

# Code Review & Security Hardener

> Eval fixture. Sandbox mock — not a live vendor API.

Audits source code for OWASP Top 10 vulnerabilities, SQL injection, XSS, buffer overflows, and architectural smells

```json
{
  "name": "code_security_audit",
  "description": "Audit source code for security vulnerabilities and architectural defects",
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
