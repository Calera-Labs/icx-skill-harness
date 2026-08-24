---
name: "Terraform Infrastructure Planner"
id: "terraform_infra_plan"
execution: "sandbox-mock"
category: "DevOps"
description: "Generates and validates Terraform execution plans for cloud infrastructure"
triggers: ["terraform","hcl","infrastructure","iac","terraform plan"]
keywords: ["terraform","iac","hcl","infrastructure","plan","apply"]
---

# Terraform Infrastructure Planner

> Eval fixture. Sandbox mock — not a live vendor API.

Generates and validates Terraform execution plans for cloud infrastructure

```json
{
  "name": "terraform_plan",
  "description": "Run terraform plan and validate HCL",
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
