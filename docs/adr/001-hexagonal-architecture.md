# 1. Use Hexagonal Architecture and Modular Monolith

## Status

Accepted

## Context

The BazaarGo system handles complex business rules (pricing, stock, authentication, and AI-driven chat). We need a way to isolate these rules from external dependencies (HTTP frameworks, databases, and LLM providers). Additionally, deploying a full microservices architecture with multiple separate databases would introduce unnecessary operational overhead for the scope of this project.

## Decision

We will use a **Modular Monolith** structure with **Hexagonal Architecture (Ports and Adapters)**.

- The core domain and use cases (Services) will reside in the center and will not depend on any infrastructure libraries.
- Communication with the outside world (Postgres, MongoDB, Gemini API, HTTP) will be handled by Adapters that implement interfaces defined in the `ports` package.
- The application will be split into two entrypoints (`cmd/api` for e-commerce and `cmd/chatbot` for RAG) to isolate their failure domains, but they will share the same codebase and infrastructure repository.

## Consequences

- **Positive:** High testability (we can mock Ports). We can swap out the LLM provider or the database without changing business logic.
- **Negative:** Requires more boilerplate (interfaces for everything). Slower initial development speed compared to a simple CRUD approach.
