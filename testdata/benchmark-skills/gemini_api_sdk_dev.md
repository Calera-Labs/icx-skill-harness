---
name: "Gemini API & SDK Developer"
id: "gemini_api_sdk_dev"
execution: "sandbox-mock"
category: "AI & Google Cloud"
description: "Develops applications using Google GenAI SDKs for Python, TypeScript, Go, function calling, and structured JSON outputs"
triggers: ["gemini sdk","google genai","structured output","function calling","gemini-flash","gemini-pro"]
keywords: ["gemini","google","genai","sdk","structured","functioncall"]
---

# Gemini API & SDK Developer

> Eval fixture. Sandbox mock — not a live vendor API.

Develops applications using Google GenAI SDKs for Python, TypeScript, Go, function calling, and structured JSON outputs

```json
{
  "name": "gemini_sdk_invoke",
  "description": "Generate code and execute Google Gemini GenAI SDK applications",
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
