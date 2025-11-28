
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

RUN go build -o search_worker ./cmd/search_worker

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/search_worker /app/search_worker
CMD ["/app/search_worker"]
