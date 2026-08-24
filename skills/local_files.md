---
name: "Local Files"
id: "local_files"
category: "Starter"
execution: "sandbox-mock"
description: "Sandbox mock of a workspace file-read skill. This harness does not read disk."
triggers: ["read file", "workspace file", "local path"]
keywords: ["file", "filesystem", "workspace", "mock"]
---

# Local Files

Default tool execution in this repo is a sandbox mock. It does not open files.

```json
{
  "name": "local_files_read",
  "description": "Sandbox mock of a workspace file read",
  "parameters": {
    "type": "object",
    "properties": {
      "relative_path": {"type": "string", "description": "Path relative to the workspace root (not read)"}
    },
    "required": ["relative_path"]
  }
}
```
