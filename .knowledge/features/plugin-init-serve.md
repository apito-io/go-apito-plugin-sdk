---
type: feature
title: Plugin Init Serve
description: Init factory, Plugin.Serve, and HashiCorp go-plugin gRPC server lifecycle
resource: sdk.go
tags: [go-plugin-build-sdk, plugin, serve, grpc]
timestamp: 2026-07-07T00:00:00Z
---

# Plugin Init Serve

## Purpose

Go entry for Apito engine plugins. `Init(name, version, apiKey)` builds `Plugin`; `Serve()` starts HashiCorp go-plugin gRPC server implementing `PluginService`. Protocol: [plugin-grpc-protocol](../../../../.knowledge/features/plugin-grpc-protocol.md).

## Flows

- **Init**: `plugin := sdk.Init("my-plugin", "1.0.0", apiKey)`.
- **Register**: queries, mutations, REST, functions on `Plugin` before serve.
- **Serve**: `plugin.Serve()` — handshake, block until shutdown.
- **Main**: `examples/main.go` pattern — `plugin.Serve()` as process entry.

## Main files

- `sdk.go` — `Init`, `Plugin`, `Serve`, gRPC service impl
- `helpers.go` — registration sugar
- `examples/main.go` — minimal plugin
- `Makefile` — build plugin binary

## Dependencies

- `github.com/hashicorp/go-plugin`
- `github.com/apito-io/types/protobuff`
- Global: [plugin-grpc-protocol](../../../../.knowledge/features/plugin-grpc-protocol.md)

## Invariants

- All registrations before `Serve()` — no hot-reload.
- Plugin binary must be executable by engine plugin manager.
- Use hclog for logs — not fmt to stdout during handshake.

## Common bugs

- `Serve()` called from multiple goroutines.
- API key mismatch — engine rejects plugin at registration.
- Wrong plugin binary path in engine config.

## Tests

- `examples/` build + engine local attach
- Compare handshake with JS [grpc-proto-handshake](../js-plugin-build-sdk/.knowledge/features/grpc-proto-handshake.md)

## Related

- [graphql-type-system](graphql-type-system.md), [health-diagnostics](health-diagnostics.md)
- `PLUGIN_DEVELOPMENT_GUIDE.md`
