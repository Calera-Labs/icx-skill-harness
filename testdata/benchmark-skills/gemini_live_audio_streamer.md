---
name: "Gemini Live API & Audio Streamer"
id: "gemini_live_audio_streamer"
execution: "sandbox-mock"
category: "AI & Google Cloud"
description: "Connects to Gemini Live API via WebSockets for real-time bidirectional audio/video streaming and voice activity detection"
triggers: ["gemini live","live api","websocket audio","bidirectional streaming","voice activity detection"]
keywords: ["live","audio","streaming","websocket","vad","voice"]
---

# Gemini Live API & Audio Streamer

> Eval fixture. Sandbox mock — not a live vendor API.

Connects to Gemini Live API via WebSockets for real-time bidirectional audio/video streaming and voice activity detection

```json
{
  "name": "gemini_live_stream",
  "description": "Manage real-time WebSocket streaming with Gemini Live API",
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
