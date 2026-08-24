---
name: "Gemini Omni Video & Multimodal Editor"
id: "gemini_omni_video_editor"
execution: "sandbox-mock"
category: "AI & Google Cloud"
description: "Performs generative video editing, frame-to-video transitions, ffmpeg preprocessing, and synchronized audio reconstruction"
triggers: ["omni flash","video editing","generative video","text to video","ffmpeg video"]
keywords: ["video","omni","multimodal","ffmpeg","transition","animation"]
---

# Gemini Omni Video & Multimodal Editor

> Eval fixture. Sandbox mock — not a live vendor API.

Performs generative video editing, frame-to-video transitions, ffmpeg preprocessing, and synchronized audio reconstruction

```json
{
  "name": "gemini_omni_edit",
  "description": "Execute generative video editing with Gemini Omni Flash",
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
