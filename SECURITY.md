# Security Policy

## Reporting a vulnerability

Use GitHub **Private vulnerability reporting** on this repository, or email **support@caleralabs.com**.

Include the affected command or endpoint, a minimal reproduction, and whether secrets were involved.

## Defaults

- The gateway binds `127.0.0.1` and requires a bearer gateway token.
- Model and ICX keys belong in environment variables, not command-line flags (`-icx-key` / `-byok-key` / `-gemini-key` are rejected).
- Model keys are not sent to ICX.
- Hosted ICX HTTP failures are returned as errors. Process-local storage is labeled `LOCAL_FALLBACK`.
- State crystallization is off unless `-crystallize` is passed.
- `POST /v1/skills/register` is disabled unless `--enable-hot-register` is set.

## Secrets

Never commit `.env` or `.icx-gateway-token`. If a key was exposed (including in a stream, screenshot, or OneDrive copy), rotate it at the issuer before any public push.
