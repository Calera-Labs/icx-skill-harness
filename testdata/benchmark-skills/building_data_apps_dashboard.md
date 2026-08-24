---
name: "Data App & Dashboard Builder"
id: "building_data_apps_dashboard"
execution: "sandbox-mock"
category: "Frontend & Apps"
description: "Scaffolds interactive data visualization web apps using React + Vite or Streamlit with Gemini Data Analytics chat"
triggers: ["data app","dashboard","streamlit","react vite","data visualization","chat with your data"]
keywords: ["dashboard","streamlit","react","frontend","charts","analytics"]
---

# Data App & Dashboard Builder

> Eval fixture. Sandbox mock — not a live vendor API.

Scaffolds interactive data visualization web apps using React + Vite or Streamlit with Gemini Data Analytics chat

```json
{
  "name": "build_data_app_scaffold",
  "description": "Scaffold React or Streamlit data applications with BigQuery integration",
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
