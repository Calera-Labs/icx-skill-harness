---
name: "RabbitMQ Message Broker Consumer"
id: "rabbitmq_queue_consumer"
execution: "sandbox-mock"
category: "Data"
description: "Acknowledges messages and manages dead-letter queues in RabbitMQ"
triggers: ["rabbitmq","amqp","queue","dead letter","message broker"]
keywords: ["rabbitmq","amqp","queue","broker","routing"]
---

# RabbitMQ Message Broker Consumer

> Eval fixture. Sandbox mock — not a live vendor API.

Acknowledges messages and manages dead-letter queues in RabbitMQ

```json
{
  "name": "rabbitmq_manage_queue",
  "description": "Manage RabbitMQ message queues",
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
