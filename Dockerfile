FROM golang:1.27.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bazaargo-api ./cmd/api/main.go

FROM alpine:3.24.4
WORKDIR /app

RUN apk --no-cache add \
    ca-certificates=20261020-r0 \
    tzdata=2026f-r0 \
    && adduser -D -g '' appuser

COPY --from=builder /app/bazaargo-api .

RUN chown -R appuser:appuser /app


USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8080/health || exit 1

CMD ["./bazaargo-api"]
