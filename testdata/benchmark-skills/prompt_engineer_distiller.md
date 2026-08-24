---
name: "Prompt Optimization & Distiller"
id: "prompt_engineer_distiller"
execution: "sandbox-mock"
category: "AI & Prompting"
description: "Refines system instructions, few-shot prompt exemplars, DSPy-style signatures, and chain-of-thought constraints"
triggers: ["prompt engineering","prompt optimization","few-shot","dspy","system prompt","distill prompt"]
keywords: ["prompt","instruction","fewshot","cot","llm","tuning"]
---

# Prompt Optimization & Distiller

> Eval fixture. Sandbox mock — not a live vendor API.

Refines system instructions, few-shot prompt exemplars, DSPy-style signatures, and chain-of-thought constraints

```json
{
  "name": "prompt_distill_optimize",
  "description": "Optimize, distill, and benchmark LLM system prompts and few-shot examples",
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
