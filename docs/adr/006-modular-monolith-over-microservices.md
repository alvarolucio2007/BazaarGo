# 6. Use a Modular Monolith over Microservices for RAG Integration

Date: 2026-11-26

## Status

Accepted (Rejection of Microservices Architecture)

## Context

Upon receiving the AI RAG implementation request, we evaluated two architectural paths: extending the existing Modular Monolith or splitting the application into a Microservices architecture.

While microservices provide greater scalability and independent maintainability, they introduce unnecessary operational complexity, network latency, and distributed transaction overhead. Given the short deadline and the scope of the project, the trade-off of adopting a microservices architecture is too high.

## Decision

We will implement the RAG feature as a new module within our Modular Monolith (with a dedicated entrypoint, `cmd/chatbot`), rather than building a fully isolated microservice. We will adapt the existing infrastructure to support it, notably by swapping the `postgres:18-alpine` image in the `docker-compose.yml` for `pgvector:pg18-trixie`, which has the pgvector extension built-in.

## Consequences

- **Positive:** Faster development cycle, easier local debugging, and slightly lower latency (no inter-service network calls). The shared infrastructure (Traefik, LGTM) is reused seamlessly.
- **Negative:** The database acts as a Single Point of Failure for the whole system. It is less scalable than independent microservices. However, because the code is strictly isolated within Hexagonal boundaries, the RAG module can be easily extracted into a dedicated microservice in the future if the need arises.
