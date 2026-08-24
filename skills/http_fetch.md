---
name: "HTTP Fetch"
id: "http_fetch"
category: "Starter"
execution: "sandbox-mock"
description: "Sandbox mock of an HTTPS GET skill. This harness does not fetch URLs."
triggers: ["http get", "fetch url", "download page"]
keywords: ["http", "https", "fetch", "url", "mock"]
---

# HTTP Fetch

Default tool execution in this repo is a sandbox mock. It does not perform network requests.

```json
{
  "name": "http_fetch",
  "description": "Sandbox mock: GET a public HTTPS URL",
  "parameters": {
    "type": "object",
    "properties": {
      "url": {"type": "string", "description": "https URL (not fetched)"}
    },
    "required": ["url"]
  }
}
```
