---
name: "LaunchDarkly Feature Flag Controller"
id: "launchdarkly_feature_flags"
execution: "sandbox-mock"
category: "DevOps"
description: "Toggles feature flags and manages rollout percentage rules"
triggers: ["launchdarkly","feature flag","toggle","rollout","kill switch"]
keywords: ["launchdarkly","feature","flag","rollout","toggle"]
---

# LaunchDarkly Feature Flag Controller

> Eval fixture. Sandbox mock — not a live vendor API.

Toggles feature flags and manages rollout percentage rules

```json
{
  "name": "launchdarkly_toggle_flag",
  "description": "Toggle LaunchDarkly feature flags",
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
