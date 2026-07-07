---
type: feature
title: REST Endpoints Multipart
description: REST endpoint registration with JSON and multipart request handling in Go plugins
resource: sdk.go
tags: [go-plugin-build-sdk, rest, multipart, endpoints]
timestamp: 2026-07-07T00:00:00Z
---

# REST Endpoints Multipart

## Purpose

Register plugin HTTP endpoints with method, path, optional schema, and `RESTHandlerFunc`. Engine proxies HTTP to plugin; handlers receive parsed path/query/body including **multipart** uploads where applicable.

## Flows

- **Register**: `RegisterRESTAPI(RESTEndpoint{Method, Path, Handler, Schema?}, handler)`.
- **Batch**: `RegisterRESTAPIs(endpoints, handlersMap)`.
- **Handler**: `RESTHandlerFunc(ctx, args)` — args include path params, query, body, files.
- **Multipart**: large uploads stream through engine → plugin gRPC args.

## Main files

- `sdk.go` — `RESTEndpoint`, `RegisterRESTAPI`, handler dispatch
- `helpers.go` — REST registration helpers
- `PLUGIN_DEVELOPMENT_GUIDE.md` — REST section

## Dependencies

- [plugin-init-serve](plugin-init-serve.md)
- [coded-and-graphql-errors](coded-and-graphql-errors.md) for HTTP status mapping

## Invariants

- Handler keys unique per plugin REST map.
- Return JSON-marshalable values or `CodedError` for status codes.
- Path patterns must not collide with engine core routes.

## Common bugs

- Assuming `args` struct shape without reading engine version docs.
- Multipart field names mismatch — empty file in handler.
- Returning Go error without `CodedError` — engine maps to 500 only.

## Tests

- Example REST handler in `examples/main.go`
- Integration test upload through engine HTTP → plugin

## Related

- JS: [rest-endpoint-registration](../js-plugin-build-sdk/.knowledge/features/rest-endpoint-registration.md)
- [custom-functions-and-context](../js-plugin-build-sdk/.knowledge/features/custom-functions-and-context.md)
