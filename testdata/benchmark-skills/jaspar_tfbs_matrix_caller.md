---
name: "JASPAR Transcription Factor Binding Profiles"
id: "jaspar_tfbs_matrix_caller"
execution: "sandbox-mock"
category: "Genomics"
description: "Queries JASPAR database for Transcription Factor binding Position Weight Matrices (PWM) and MEME formatted motifs"
triggers: ["jaspar","transcription factor","tfbs","pwm matrix","binding motif","meme format"]
keywords: ["jaspar","transcription","tfbs","pwm","motif","dna binding"]
---

# JASPAR Transcription Factor Binding Profiles

> Eval fixture. Sandbox mock — not a live vendor API.

Queries JASPAR database for Transcription Factor binding Position Weight Matrices (PWM) and MEME formatted motifs

```json
{
  "name": "jaspar_matrix_query",
  "description": "Fetch transcription factor binding profiles and PWM matrices from JASPAR",
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
