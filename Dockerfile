# Multi-stage Dockerfile for Kid Dictionary
# - dev: Development with air hot reload
# - prod: Optimized production binary

# =============================================================================
# Base stage: Common Go setup
# =============================================================================
FROM golang:1.25-alpine AS base

WORKDIR /app

# Install required tools
RUN apk add --no-cache git

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# =============================================================================
# Development stage: Air hot reload
# =============================================================================
FROM base AS dev

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy air configuration
COPY .air.toml ./

# Copy source code (will be overridden by volume mount in docker-compose)
COPY . .

# Expose port
EXPOSE 8080

# Run air for hot reload
CMD ["air", "-c", ".air.toml"]

# =============================================================================
# Builder stage: Compile optimized binary
# =============================================================================
FROM base AS builder

# Copy source code
COPY . .

# Build optimized binary with CGO disabled
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# =============================================================================
# Production stage: Minimal runtime image
# =============================================================================
FROM alpine:3.19 AS prod

# Add CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the server
CMD ["./server"]
