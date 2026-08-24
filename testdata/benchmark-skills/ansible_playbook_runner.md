---
name: "Ansible Automation Playbook Runner"
id: "ansible_playbook_runner"
execution: "sandbox-mock"
category: "DevOps"
description: "Executes configuration management playbooks across fleet of servers"
triggers: ["ansible","playbook","automation","sysadmin","inventory"]
keywords: ["ansible","playbook","yaml","ssh","automation"]
---

# Ansible Automation Playbook Runner

> Eval fixture. Sandbox mock — not a live vendor API.

Executes configuration management playbooks across fleet of servers

```json
{
  "name": "ansible_run_playbook",
  "description": "Execute Ansible configuration playbooks",
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
