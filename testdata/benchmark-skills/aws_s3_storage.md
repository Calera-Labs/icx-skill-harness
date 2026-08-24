---
name: "AWS S3 Cloud Storage Manager"
id: "aws_s3_storage"
execution: "sandbox-mock"
category: "Cloud"
description: "Uploads, downloads, and manages bucket lifecycle policies on AWS S3"
triggers: ["s3","aws","bucket","storage","aws s3"]
keywords: ["s3","aws","bucket","object","storage"]
---

# AWS S3 Cloud Storage Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Uploads, downloads, and manages bucket lifecycle policies on AWS S3

```json
{
  "name": "s3_bucket_ops",
  "description": "Manage AWS S3 storage buckets and objects",
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
