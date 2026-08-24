---
name: "Cloud Storage Security & SAIF Auditor"
id: "gcs_security_saif_auditor"
execution: "sandbox-mock"
category: "Security & Google Cloud"
description: "Assesses Google Cloud Storage bucket IAM permissions, public access prevention, SAIF compliance, and KMS CMEK keys"
triggers: ["gcs security","saif compliance","bucket iam","public access prevention","cmek","gcs audit"]
keywords: ["gcs","security","saif","cloud storage","bucket","iam"]
---

# Cloud Storage Security & SAIF Auditor

> Eval fixture. Sandbox mock — not a live vendor API.

Assesses Google Cloud Storage bucket IAM permissions, public access prevention, SAIF compliance, and KMS CMEK keys

```json
{
  "name": "gcs_security_audit",
  "description": "Audit Google Cloud Storage buckets for SAIF compliance and security risks",
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
