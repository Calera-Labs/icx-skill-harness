---
name: "Design System & Theme Factory"
id: "theme_factory_stylist"
execution: "sandbox-mock"
category: "Design & UI"
description: "Generates cohesive UI color palettes, design tokens, typography pairings, and WCAG AA contrast matrices"
triggers: ["theme","palette","design token","typography","wcag","color scheme","ui design"]
keywords: ["theme","colors","tokens","styling","css","dark mode","tailwind"]
---

# Design System & Theme Factory

> Eval fixture. Sandbox mock — not a live vendor API.

Generates cohesive UI color palettes, design tokens, typography pairings, and WCAG AA contrast matrices

```json
{
  "name": "theme_factory_generate",
  "description": "Generate design system color tokens and typography palettes",
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
