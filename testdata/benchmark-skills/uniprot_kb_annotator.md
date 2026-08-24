---
name: "UniProt Knowledgebase Annotator"
id: "uniprot_kb_annotator"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Fetches UniProtKB protein sequence records, active catalytic residues, isoforms, subcellular locations, and GO tags"
triggers: ["uniprot","uniprotkb","protein sequence","catalytic site","isoform","protein annotation"]
keywords: ["uniprot","protein","sequence","fasta","annotation","swissprot"]
---

# UniProt Knowledgebase Annotator

> Eval fixture. Sandbox mock — not a live vendor API.

Fetches UniProtKB protein sequence records, active catalytic residues, isoforms, subcellular locations, and GO tags

```json
{
  "name": "uniprot_fetch_protein",
  "description": "Look up protein metadata and sequences in UniProtKB",
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
