# 2. Contract-First API Design using OpenAPI

## Status

Accepted

## Context

In previous projects, defining HTTP handlers first and generating documentation later led to documentation drift and mismatched types between the backend and frontend. We need a strict, single source of truth for our API contracts.

## Decision

We will use **Spec-First development** with OpenAPI 3.0. The `openapi.yaml` file is the source of truth. We use `oapi-codegen` to generate the Go types, interfaces, and the Strict Server boilerplate.

## Consequences

- **Positive:** Compile-time safety. If the code doesn't match the contract, it doesn't compile. The Swagger UI is automatically generated and always up-to-date. Frontend mocking becomes trivial.
- **Negative:** Writing YAML can be verbose. Modifying an endpoint requires re-running the generation command, which adds a step to the development loop.
