---
type: feature
title: Coded and GraphQL Errors
description: CodedError with HTTP status and GraphQLError with extensions for plugin handlers
resource: sdk.go
tags: [go-plugin-build-sdk, errors, graphql, http]
timestamp: 2026-07-07T00:00:00Z
---

# Coded and GraphQL Errors

## Purpose

Unified error types for REST (`CodedError` + HTTP status) and GraphQL (`GraphQLError` + extensions). Constructor helpers mirror JS plugin SDK error patterns.

## Flows

- **HTTP**: `BadRequestError`, `UnauthorizedError`, `NotFoundError`, `InternalServerError`, etc.
- **GraphQL**: `GraphQLErrorWithCode`, `GraphQLErrorWithExtensions`, location/path fields.
- **REST handler**: return `error` implementing `CodedError` → engine sets HTTP status.
- **GraphQL resolver**: return `GraphQLError` → engine formats `errors[]` response.

## Main files

- `sdk.go` — `CodedError`, `GraphQLError`, constructors (top of file)
- `helpers.go` — additional error utilities if present

## Dependencies

- [graphql-type-system](graphql-type-system.md) resolvers
- [rest-endpoints-multipart](rest-endpoints-multipart.md) handlers

## Invariants

- Prefer typed error constructors over `fmt.Errorf` in plugin boundaries.
- GraphQL extensions `code` field should use engine-recognized codes.
- HTTP 4xx for client mistakes, 5xx only for true internal failures.

## Common bugs

- Returning wrapped `CodedError` — engine may not unwrap — return direct type.
- GraphQL error without message — empty client toast.
- Mixing HTTP errors in GraphQL resolver — use GraphQL error type instead.

## Tests

- Handler returns `BadRequestError` → HTTP 400 from engine proxy
- Resolver returns `GraphQLErrorWithCode` → extensions.code in GraphQL JSON

## Related

- JS: [error-handling-graphql](../js-plugin-build-sdk/.knowledge/features/error-handling-graphql.md)
- `PLUGIN_DEVELOPMENT_GUIDE.md`
