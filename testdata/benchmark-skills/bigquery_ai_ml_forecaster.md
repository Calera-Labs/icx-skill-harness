---
name: "BigQuery ML & Time Series Forecaster"
id: "bigquery_ai_ml_forecaster"
execution: "sandbox-mock"
category: "Data & Google Cloud"
description: "Builds BigQuery ML ARIMA_PLUS models, autoencoder anomaly detection, and integrates remote Gemini LLM models in SQL"
triggers: ["bigquery ml","arima_plus","time series forecast","bqml","anomaly detection"]
keywords: ["bqml","ml","forecast","arima","bigquery","machine learning"]
---

# BigQuery ML & Time Series Forecaster

> Eval fixture. Sandbox mock — not a live vendor API.

Builds BigQuery ML ARIMA_PLUS models, autoencoder anomaly detection, and integrates remote Gemini LLM models in SQL

```json
{
  "name": "bigquery_ml_forecast",
  "description": "Train and evaluate BigQuery ML models and time-series forecasts",
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
