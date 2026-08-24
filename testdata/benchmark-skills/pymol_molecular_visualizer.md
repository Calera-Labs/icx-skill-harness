---
name: "PyMOL 3D Molecular Rendering Engine"
id: "pymol_molecular_visualizer"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Generates PyMOL rendering scripts, superpositions, binding pocket raytracing, and surface electrostatic views"
triggers: ["pymol","molecular visualization","render protein","binding pocket render","structural alignment"]
keywords: ["pymol","rendering","3d","visualization","raytracing","superposition"]
---

# PyMOL 3D Molecular Rendering Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Generates PyMOL rendering scripts, superpositions, binding pocket raytracing, and surface electrostatic views

```json
{
  "name": "pymol_render_scene",
  "description": "Generate PyMOL visualization scripts and render molecular scenes",
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
