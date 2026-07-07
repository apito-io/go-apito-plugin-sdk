---
type: feature
title: Complex Types and Parsing
description: Complex GraphQL types, protobuf Any/Struct parsing, and reflection helpers in Go plugins
resource: sdk.go
tags: [go-plugin-build-sdk, types, protobuf, parsing]
timestamp: 2026-07-07T00:00:00Z
---

# Complex Types and Parsing

## Purpose

Handle nested GraphQL inputs, maps, lists, and protobuf `anypb`/`structpb` values crossing the gRPC plugin boundary. Documented examples in `COMPLEX_TYPES_EXAMPLES.md`.

## Flows

- **Define**: nested `GraphQLTypeDefinition` trees for input objects.
- **Parse args**: reflection + `structpb` conversion from gRPC request payload.
- **Return**: marshal complex Go structs to GraphQL-compatible JSON values.
- **Any types**: `google.golang.org/protobuf/types/known/anypb` for dynamic payloads.

## Main files

- `sdk.go` — type parsing, `structpb`, `anypb` usage in dispatch
- `COMPLEX_TYPES_EXAMPLES.md` — worked examples
- `TYPE_SYSTEM.md` — type composition rules
- `helpers.go` — arg extraction helpers

## Dependencies

- [graphql-type-system](graphql-type-system.md)
- `google.golang.org/protobuf`
- Engine plugin RPC payload format

## Invariants

- JSON numbers arrive as float64 in maps — cast carefully to int fields.
- Nil pointer vs empty slice semantics matter for GraphQL list fields.
- Do not mutate `structpb` values shared across handler invocations.

## Common bugs

- `map[string]interface{}` nested key type assertion panics — validate first.
- Protobuf Any unpack wrong message type.
- GraphQL input optional field omitted vs null confusion.

## Tests

- `COMPLEX_TYPES_EXAMPLES.md` compile + run against engine
- Round-trip complex mutation input in integration test

## Related

- [graphql-type-system](graphql-type-system.md)
- `COMPLEX_TYPES_EXAMPLES.md`, `TYPE_SYSTEM.md`
