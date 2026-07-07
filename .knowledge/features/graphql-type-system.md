---
type: feature
title: GraphQL Type System
description: GraphQL field and object type registration with Go type definitions for plugins
resource: sdk.go
tags: [go-plugin-build-sdk, graphql, types, registration]
timestamp: 2026-07-07T00:00:00Z
---

# GraphQL Type System

## Purpose

Register plugin GraphQL queries/mutations with typed field definitions. Supports nested `GraphQLTypeDefinition`, object type registry, and resolver function binding. See also `TYPE_SYSTEM.md` for extended examples.

## Flows

- **Single op**: `RegisterQuery(name, GraphQLField{…}, resolver)`.
- **Batch**: `RegisterQueries(fieldsMap, resolversMap)`.
- **Object types**: register via object type map before referencing in fields.
- **Resolvers**: `ResolverFunc(context, args) (any, error)`.

## Main files

- `sdk.go` — `GraphQLField`, `GraphQLTypeDefinition`, register methods
- `TYPE_SYSTEM.md` — type system documentation
- `COMPLEX_TYPES_EXAMPLES.md` — advanced type patterns
- `helpers.go` — builder helpers

## Dependencies

- [plugin-init-serve](plugin-init-serve.md)
- [complex-types-and-parsing](complex-types-and-parsing.md) for nested types
- [coded-and-graphql-errors](coded-and-graphql-errors.md) in resolvers

## Invariants

- Resolver map keys must match field `Resolve` name strings.
- List/NonNull wrappers use nested `GraphQLTypeDefinition.OfType`.
- Object types must be registered before fields reference them.

## Common bugs

- `interface{}` return type mismatch with declared GraphQL field type.
- Circular object type references without forward registration.
- Mutation resolver returning pointer vs value inconsistently.

## Tests

- Build example plugin; verify merged engine schema
- `TYPE_SYSTEM.md` examples compile

## Related

- JS: [graphql-registration-helpers](../js-plugin-build-sdk/.knowledge/features/graphql-registration-helpers.md)
- `TYPE_SYSTEM.md`, `COMPLEX_TYPES_EXAMPLES.md`
