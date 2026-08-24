---
name: "ClinicalTrials.gov Protocol Matcher"
id: "clinical_trials_gov_api"
execution: "sandbox-mock"
category: "Life Sciences"
description: "Queries ClinicalTrials.gov APIv2 for interventional trials, disease conditions, NCT identifiers, and eligibility criteria"
triggers: ["clinical trials","clinicaltrials.gov","nct","trial eligibility","phase 3 trial","interventional study"]
keywords: ["trials","nct","clinical","pharma","recruiting","study"]
---

# ClinicalTrials.gov Protocol Matcher

> Eval fixture. Sandbox mock — not a live vendor API.

Queries ClinicalTrials.gov APIv2 for interventional trials, disease conditions, NCT identifiers, and eligibility criteria

```json
{
  "name": "clinical_trials_lookup",
  "description": "Search ClinicalTrials.gov for clinical study protocols and patient matching",
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
