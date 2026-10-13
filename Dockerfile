#############################
# Builder
#############################
FROM golang:1.25-alpine3.21 AS builder

ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOTOOLCHAIN=auto

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Copy dependency terlebih dahulu
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API
RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o engine \
    ./cmd/api

#############################
# Runtime
#############################
FROM alpine:3.21

ENV TZ=Asia/Jakarta

RUN apk add --no-cache \
    ca-certificates \
    tzdata && \
    addgroup -S app && \
    adduser -S app -G app

WORKDIR /app

# Binary
COPY --from=builder /app/engine /app/engine

# Config
COPY config.yml /app/config.yml

RUN chown -R app:app /app && \
    chmod 500 /app/engine && \
    chmod 644 /app/config.yml

USER app

EXPOSE 8080

ENTRYPOINT ["/app/engine"]
