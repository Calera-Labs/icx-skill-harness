---
name: "dbSNP Variant & rsID Coordinate Mapper"
id: "dbsnp_variant_mapper"
execution: "sandbox-mock"
category: "Genomics"
description: "Maps rsIDs to GRCh38 genomic coordinates, HGVS strings, minor allele frequencies, and indel classifications"
triggers: ["dbsnp","rsid","snp","genomic coordinates","grch38 variant","indel lookup"]
keywords: ["dbsnp","rsid","snp","genomics","allele","ncbi"]
---

# dbSNP Variant & rsID Coordinate Mapper

> Eval fixture. Sandbox mock — not a live vendor API.

Maps rsIDs to GRCh38 genomic coordinates, HGVS strings, minor allele frequencies, and indel classifications

```json
{
  "name": "dbsnp_lookup_variant",
  "description": "Resolve and map genetic variants (SNPs/indels) in NCBI dbSNP",
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
