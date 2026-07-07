---
type: feature
title: Health Diagnostics
description: Plugin health checks, hclog diagnostics, and runtime introspection for engine ops
resource: sdk.go
tags: [go-plugin-build-sdk, health, diagnostics, logging]
timestamp: 2026-07-07T00:00:00Z
---

# Health Diagnostics

## Purpose

Expose plugin health to the engine for readiness probes, debugging failed handshakes, and operational logging via `hclog`. Complements built-in function hooks similar to JS `health_check`.

## Flows

- **Health function**: register via `RegisterFunction` — engine polls plugin alive.
- **Logging**: `hclog` logger with plugin name context — avoid stdout corruption.
- **Diagnostics**: version, Go runtime, registered op counts for support bundles.
- **Failure**: Serve errors surface in engine plugin manager UI/logs.

## Main files

- `sdk.go` — health-related function registration, Serve logging
- `version.go` — plugin SDK version string
- `examples/main.go` — health registration pattern
- `PLUGIN_DEVELOPMENT_GUIDE.md` — ops section

## Dependencies

- [plugin-init-serve](plugin-init-serve.md)
- `github.com/hashicorp/go-hclog`

## Invariants

- Health check must return quickly — no blocking IO.
- Use leveled logs — Info for lifecycle, Error for handler failures.
- Do not log secrets (apiKey, user tokens) in diagnostics.

## Common bugs

- Health check performs DB ping — times out engine probe.
- fmt.Println during Serve handshake — breaks go-plugin protocol.
- Missing version in logs — hard to match plugin binary to SDK release.

## Tests

- Engine plugin list shows healthy after `Serve()`
- Kill plugin process — engine marks unhealthy

## Related

- JS: [plugin-init-serve-lifecycle](../js-plugin-build-sdk/.knowledge/features/plugin-init-serve-lifecycle.md)
- [plugin-init-serve](plugin-init-serve.md)
- Global: [plugin-grpc-protocol](../../../../.knowledge/features/plugin-grpc-protocol.md)
