# 3. Polyglot Persistence for Optimal Data Modeling

## Status

Accepted

## Context

BazaarGo handles different types of data with varying access patterns:

1. User credentials require strict ACID compliance and relational integrity.
2. The product catalog requires a flexible schema because different items have different attributes (e.g., size, color, RAM).
3. High-speed read operations (like recently viewed products) require in-memory caching.
4. AI embeddings for RAG require vector storage and similarity search.

## Decision

We will adopt a **Polyglot Persistence** approach:

- **PostgreSQL:** Used for authentication and user management. We also enabled the `pgvector` extension to store document embeddings and perform similarity searches for the RAG pipeline, avoiding the need for a dedicated vector database.
- **MongoDB:** Used for the product catalog due to its document-oriented nature.
- **Redis (Valkey):** Used for Cache-Aside implementations and tracking recently viewed items using Lists (`LPUSH`, `LTRIM`).

## Consequences

- **Positive:** The system uses the right tool for the right job, optimizing performance and schema flexibility.
- **Negative:** Increased infrastructure complexity. We need to manage multiple database connections, backups, and monitoring. Distributed transactions (2PC) are avoided by design.
