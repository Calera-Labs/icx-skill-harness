---
name: "MCP Server & Tool Manifest Builder"
id: "mcp_server_builder"
execution: "sandbox-mock"
category: "AI & MCP"
description: "Scaffolds Model Context Protocol (MCP) servers, tool declarations, JSON-RPC 2.0 schemas, and STDIO/SSE handlers"
triggers: ["mcp","model context protocol","mcp server","json-rpc","tool manifest","mcp tool"]
keywords: ["mcp","protocol","context","server","tools","manifest"]
---

# MCP Server & Tool Manifest Builder

> Eval fixture. Sandbox mock — not a live vendor API.

Scaffolds Model Context Protocol (MCP) servers, tool declarations, JSON-RPC 2.0 schemas, and STDIO/SSE handlers

```json
{
  "name": "mcp_scaffold_server",
  "description": "Scaffold and validate MCP Model Context Protocol servers and schemas",
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
