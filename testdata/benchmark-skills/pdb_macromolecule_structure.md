---
name: "Protein Data Bank (PDB) Macromolecule Structure"
id: "pdb_macromolecule_structure"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Downloads experimentally determined 3D atomic coordinates (.cif/.pdb) and binding ligand metadata from the Protein Data Bank"
triggers: ["pdb","protein data bank","pdb structure","crystal structure","pdb id","cryo-em structure"]
keywords: ["pdb","structure","crystal","macromolecule","ligand","cif"]
---

# Protein Data Bank (PDB) Macromolecule Structure

> Eval fixture. Sandbox mock — not a live vendor API.

Downloads experimentally determined 3D atomic coordinates (.cif/.pdb) and binding ligand metadata from the Protein Data Bank

```json
{
  "name": "pdb_fetch_structure",
  "description": "Fetch 3D crystal structures and experimental metadata from PDB",
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
