
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

RUN go build -o finance_service ./cmd/finance_service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/finance_service /app/finance_service
CMD ["/app/finance_service"]
