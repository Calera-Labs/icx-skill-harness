---
name: "Gemini Interactions API Agent"
id: "gemini_interactions_agent"
execution: "sandbox-mock"
category: "AI & Google Cloud"
description: "Orchestrates multi-turn conversational agents, background research loops, and multimodal interaction sessions"
triggers: ["interactions api","gemini agent","multi-turn","background research","chat session"]
keywords: ["interactions","agent","multimodal","streaming","turns"]
---

# Gemini Interactions API Agent

> Eval fixture. Sandbox mock — not a live vendor API.

Orchestrates multi-turn conversational agents, background research loops, and multimodal interaction sessions

```json
{
  "name": "gemini_interactions_call",
  "description": "Execute multi-turn conversational turns with Gemini Interactions API",
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
