---
name: "FastAPI OpenAPI Endpoint Tester"
id: "fastapi_openapi_tester"
execution: "sandbox-mock"
category: "API"
description: "Tests REST endpoints against auto-generated Swagger OpenAPI specs"
triggers: ["fastapi","openapi","swagger","rest api","endpoint test"]
keywords: ["fastapi","openapi","swagger","rest","api"]
---

# FastAPI OpenAPI Endpoint Tester

> Eval fixture. Sandbox mock — not a live vendor API.

Tests REST endpoints against auto-generated Swagger OpenAPI specs

```json
{
  "name": "fastapi_test_route",
  "description": "Test FastAPI routes and OpenAPI schemas",
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
