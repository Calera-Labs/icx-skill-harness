# Benchmark skill catalog

Eval schema cards for `-cmd diagnostic`, `-cmd benchmark`, and `-cmd e2e`.

They are **routing fixtures**, not live vendor integrations. Filenames that mention Stripe, Bloomberg, Vault, Kubernetes, and similar products are nominative labels for the mock catalog. The in-memory source is `skills.PopulateCatalog` in `pkg/skills/catalog.go`. Every file is `execution: sandbox-mock`.
