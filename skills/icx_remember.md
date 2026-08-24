---
name: "ICX Remember"
id: "icx_remember"
category: "Starter"
execution: "sandbox-mock"
description: "Sandbox mock of an ICX ingest skill. Use -cmd sync with ICX_API_KEY for a real ingest."
triggers: ["remember", "store in icx", "crystallize", "save to lattice"]
keywords: ["icx", "remember", "memory", "ingest", "mock"]
---

# ICX Remember

The default tool executor returns a fixture. Hosted ingest is `-cmd sync` (or the ICX client), not this mock tool.

```json
{
  "name": "icx_remember",
  "description": "Sandbox mock of ingesting text into an ICX space",
  "parameters": {
    "type": "object",
    "properties": {
      "text": {"type": "string", "description": "Content to store"},
      "family": {"type": "string", "description": "Optional register family"}
    },
    "required": ["text"]
  }
}
```
