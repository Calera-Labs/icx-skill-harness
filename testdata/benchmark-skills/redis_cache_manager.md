---
name: "Redis In-Memory Cache Manager"
id: "redis_cache_manager"
execution: "sandbox-mock"
category: "Database"
description: "Manages Redis keys, sets TTL, flushes cache keys, and inspects cache eviction"
triggers: ["redis","cache","ttl","flush","orderbook","redis key"]
keywords: ["redis","cache","ttl","in-memory","key-value"]
---

# Redis In-Memory Cache Manager

> Eval fixture. Sandbox mock — not a live vendor API.

Manages Redis keys, sets TTL, flushes cache keys, and inspects cache eviction

```json
{
  "name": "redis_cache_ops",
  "description": "Execute operations on Redis in-memory cache",
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
