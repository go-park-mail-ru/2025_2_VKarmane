# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .


ARG SERVICE_PATH=cmd/api
ARG BINARY_NAME=app

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/${BINARY_NAME} ./${SERVICE_PATH}


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/* /app/

COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/app"]


# docker build \
#   --build-arg SERVICE_PATH=cmd/api \
#   --build-arg BINARY_NAME=api \
#   -t api .

# docker build \
#   --build-arg SERVICE_PATH=cmd/auth_service \
#   --build-arg BINARY_NAME=auth_service \
#   -t auth_service .

# docker build \
#   --build-arg SERVICE_PATH=cmd/budget_service \
#   --build-arg BINARY_NAME=budget_service \
#   -t budget_service .


