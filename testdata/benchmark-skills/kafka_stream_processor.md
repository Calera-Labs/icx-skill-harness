---
name: "Apache Kafka Stream Processor"
id: "kafka_stream_processor"
execution: "sandbox-mock"
category: "Data"
description: "Consumes and produces high-throughput event topics across Kafka brokers"
triggers: ["kafka","stream","topic","consumer","producer","kafka broker"]
keywords: ["kafka","stream","topic","event","broker"]
---

# Apache Kafka Stream Processor

> Eval fixture. Sandbox mock — not a live vendor API.

Consumes and produces high-throughput event topics across Kafka brokers

```json
{
  "name": "kafka_topic_manager",
  "description": "Manage Kafka topics and event streams",
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
