---
name: "gRPC Protobuf Service Caller"
id: "grpc_protobuf_caller"
execution: "sandbox-mock"
category: "API"
description: "Encodes protobuf payloads and makes high-speed gRPC RPC calls"
triggers: ["grpc","protobuf","rpc","proto","grpc call"]
keywords: ["grpc","protobuf","rpc","proto","service"]
---

# gRPC Protobuf Service Caller

> Eval fixture. Sandbox mock — not a live vendor API.

Encodes protobuf payloads and makes high-speed gRPC RPC calls

```json
{
  "name": "grpc_rpc_invoke",
  "description": "Invoke gRPC service methods with protobuf",
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
