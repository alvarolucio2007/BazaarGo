FROM golang:1.27.3-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bazaargo-api ./cmd/api/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/bazaargo-api .
EXPOSE 8080
CMD ["./bazaargo-api"]
