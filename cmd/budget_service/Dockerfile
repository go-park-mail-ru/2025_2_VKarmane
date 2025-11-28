
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY migrations/ ./migrations/
COPY docker/ ./docker/

RUN go build -o budget_service ./cmd/budget_service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/budget_service /app/budget_service
CMD ["/app/budget_service"]
