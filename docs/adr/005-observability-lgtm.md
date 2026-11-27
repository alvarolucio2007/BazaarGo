# 5. OpenTelemetry and LGTM Stack for Observability

## Status

Accepted

## Context

To debug performance issues (like the latency of LLM inference) and ensure system reliability under high load, we need full visibility into our services. We need traces, metrics, and logs that can be correlated.

## Decision

We will instrument the Go APIs using the **OpenTelemetry SDK**.
The telemetry data will be exported via OTLP to the **LGTM stack** (Loki, Grafana, Tempo, Prometheus/Mimir), running via the `grafana/otel-lgtm` Docker image. We will also use **Pyroscope** for continuous profiling.

## Consequences

- **Positive:** Deep visibility. We can trace an HTTP request from the Handler down to the database query and the external LLM call. Logs are automatically enriched with `trace_id`.
- **Negative:** The observability stack is resource-intensive (requires ~1.5GB of RAM), which is a significant overhead for local development.
