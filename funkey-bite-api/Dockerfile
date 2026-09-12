# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Cache module downloads separately from source changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Mirrors `make build-local`: builds only the non-Lambda entry point.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath -ldflags="-s -w" \
    -o /out/funkey-bite-api \
    ./cmd/api/main_local.go

# ---- Runtime stage ----
FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates curl \
    && addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/funkey-bite-api ./funkey-bite-api

ENV ENVIRONMENT=production \
    PORT=8080

EXPOSE 8080
USER app

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f "http://localhost:${PORT}/health" || exit 1

ENTRYPOINT ["./funkey-bite-api"]
