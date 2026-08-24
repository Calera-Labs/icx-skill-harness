---
name: "ArXiv Scientific Preprints Search"
id: "arxiv_paper_search_fetcher"
execution: "sandbox-mock"
category: "Research & Literature"
description: "Searches arXiv repository across Computer Science, Quantitative Biology, Physics, and Math, extracting full text"
triggers: ["arxiv","preprint","arxiv id","research paper","arxiv search","cs paper"]
keywords: ["arxiv","paper","research","academic","preprint","science"]
---

# ArXiv Scientific Preprints Search

> Eval fixture. Sandbox mock — not a live vendor API.

Searches arXiv repository across Computer Science, Quantitative Biology, Physics, and Math, extracting full text

```json
{
  "name": "arxiv_search_papers",
  "description": "Search arXiv preprints and download paper metadata and text",
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
