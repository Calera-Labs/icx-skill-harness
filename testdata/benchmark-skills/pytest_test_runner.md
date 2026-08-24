---
name: "PyTest Automated Test Suite Runner"
id: "pytest_test_runner"
execution: "sandbox-mock"
category: "Testing"
description: "Executes unit and integration test fixtures and parses coverage reports"
triggers: ["pytest","test runner","coverage","fixtures","python test"]
keywords: ["pytest","test","coverage","python","unit"]
---

# PyTest Automated Test Suite Runner

> Eval fixture. Sandbox mock — not a live vendor API.

Executes unit and integration test fixtures and parses coverage reports

```json
{
  "name": "pytest_run_suite",
  "description": "Run PyTest test suite with coverage",
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
