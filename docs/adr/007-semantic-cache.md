# 7. Semantic Cache for RAG (LLM Optimization)

## Status

Accepted

## Context

Calling the Google Gemini API takes between 1.5s to 2.5s and incurs token costs. Traditional key-value caching (exact string matching) is ineffective for natural language because different users ask the same question using different words.

## Decision

We will implement a **Semantic Cache** using the existing `pgvector` database.
Before calling the LLM, the incoming question is embedded and compared against a `semantic_cache` table using cosine similarity. If the similarity is >= 0.95, the stored answer is returned immediately. If not (Cache Miss), the RAG flow executes, and the new question/answer pair is saved to the cache.

## Consequences

- **Positive:** Latency for repeated/similar questions drops from ~2.5s to ~350ms (a 7x improvement). Saves API costs and provides resilience against LLM provider outages (e.g., 503 errors).
- **Negative:** Even on a Cache Hit, we still have to generate the embedding for the user's question (costing ~350ms). False positives could occur if the similarity threshold is set too low.
