---
name: "Dataproc Serverless & Apache Spark"
id: "gcp_spark_dataproc_serverless"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Submits PySpark and Spark SQL batch jobs, manages Dataproc Serverless sessions, and queries BigLake Iceberg tables"
triggers: ["dataproc","spark","pyspark","spark sql","iceberg catalog","dataproc serverless"]
keywords: ["spark","dataproc","pyspark","iceberg","biglake","serverless"]
---

# Dataproc Serverless & Apache Spark

> Eval fixture. Sandbox mock — not a live vendor API.

Submits PySpark and Spark SQL batch jobs, manages Dataproc Serverless sessions, and queries BigLake Iceberg tables

```json
{
  "name": "dataproc_spark_submit",
  "description": "Submit and monitor PySpark batch jobs on Dataproc Serverless",
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
