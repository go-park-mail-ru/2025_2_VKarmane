# ---------- Build stage ----------
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

ARG SERVICE_PATH=cmd/api
ARG BINARY_NAME=app

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/${BINARY_NAME} ./${SERVICE_PATH}

# ---------- Runtime stage ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/* /app/
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/app"]
