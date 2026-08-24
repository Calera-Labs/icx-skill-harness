# Contributing

## Setup

```bash
go test ./...
go vet ./...
```

Do not commit secrets. Use `.env.example` as the template. Do not pass API keys on the command line.

## Pull requests

- Keep changes scoped to the router, parser, gateway, BYOK adapters, tests, or docs.
- Docs should match the binary.
- Gateway bind stays loopback by default.
- Do not present eval fixtures as live vendor integrations.

Public skills live in `/skills` (sandbox-mock starters). Eval fixtures live in `testdata/benchmark-skills`.
