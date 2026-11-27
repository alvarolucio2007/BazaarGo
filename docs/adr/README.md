# Architecture Decision Records (ADRs)

This directory contains the Architecture Decision Records (ADRs) for the **BazaarGo** project.

## What is an ADR?

An Architecture Decision Record is a short text file that captures an important architectural decision made along with its context and consequences. We use this to track the "why" behind our technical choices, ensuring that future maintainers understand the constraints and trade-offs at the time of the decision.

## List of Decisions

| ID                                                  | Title                                           | Status   |
| :-------------------------------------------------- | :---------------------------------------------- | :------- |
| [001](./001-hexagonal-architecture.md)              | Use Hexagonal Architecture and Modular Monolith | Accepted |
| [002](./002-openapi-spec-first.md)                  | Contract-First API Design using OpenAPI         | Accepted |
| [003](./003-polyglot-persistence.md)                | Polyglot Persistence for Optimal Data Modeling  | Accepted |
| [004](./004-paseto-over-jwt.md)                     | Use PASETO v4 instead of JWT for Authentication | Accepted |
| [005](./005-observability-lgtm.md)                  | OpenTelemetry and LGTM Stack for Observability  | Accepted |
| [006](./006-modular-monolith-over-microservices.md) | Use Modular Monolith over Microservices for RAG | Accepted |
| [007](./007-semantic-cache.md)                      | Semantic Cache for RAG (LLM Optimization)       | Accepted |

---

_Note: The file names above assume the naming convention `XXX-title.md`. Make sure your ADR files match these names._
