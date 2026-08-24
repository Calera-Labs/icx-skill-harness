---
name: "AWS Lambda Serverless Function Manager"
id: "aws_lambda_serverless"
execution: "sandbox-mock"
category: "Cloud"
description: "Deploys Python/Go functions and configures EventBridge trigger rules"
triggers: ["aws lambda","serverless","eventbridge","lambda function"]
keywords: ["lambda","aws","serverless","eventbridge","functions"]
---

# AWS Lambda Serverless Function Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Deploys Python/Go functions and configures EventBridge trigger rules

```json
{
  "name": "aws_lambda_invoke",
  "description": "Invoke and deploy AWS Lambda serverless functions",
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
