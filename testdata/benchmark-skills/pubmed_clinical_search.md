---
name: "PubMed Biomedical Literature Search"
id: "pubmed_clinical_search"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Searches PubMed NCBI biomedical publications, fetches abstracts, clinical trial citations, and MeSH indexing"
triggers: ["pubmed","ncbi","pmid","biomedical literature","clinical trial paper","mesh terms"]
keywords: ["pubmed","literature","journal","mesh","medicine","biology"]
---

# PubMed Biomedical Literature Search

> Eval fixture. Sandbox mock — not a live vendor API.

Searches PubMed NCBI biomedical publications, fetches abstracts, clinical trial citations, and MeSH indexing

```json
{
  "name": "pubmed_search_articles",
  "description": "Search PubMed database for peer-reviewed biomedical literature",
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
