---
name: "Neo4j Graph Database & Cypher Engine"
id: "neo4j_graph_cypher"
execution: "sandbox-mock"
category: "Database"
description: "Executes Cypher graph queries to find shortest paths and entity clusters"
triggers: ["neo4j","cypher","graph","nodes","relationships","graph traversal"]
keywords: ["neo4j","cypher","graph","nodes","relationships"]
---

# Neo4j Graph Database & Cypher Engine

> Eval fixture. Sandbox mock — not a live vendor API.

Executes Cypher graph queries to find shortest paths and entity clusters

```json
{
  "name": "neo4j_cypher_query",
  "description": "Run Cypher queries on Neo4j graph database",
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
