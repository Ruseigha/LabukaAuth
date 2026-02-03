FROM golang:1.21-alpine AS builder

# Install build dependencies
# WHY: Need git for go mod download, ca-certificates for HTTPS
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files first
# WHY: Docker layer caching - dependencies change less frequently
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
# WHY: CGO_ENABLED=0 creates static binary (no C dependencies)
# -ldflags="-w -s" strips debug info (smaller binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o auth-service \
    cmd/server/main.go

# Stage 2: Runtime
# WHY: Minimal image for running the binary
FROM alpine:latest

# Install runtime dependencies
# WHY: ca-certificates for HTTPS, tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
# WHY: Security - don't run as root
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/auth-service .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose ports
EXPOSE 8001 9001

# Health check
# WHY: Kubernetes/Docker can check if container is healthy
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8001/health || exit 1

# Run the binary
ENTRYPOINT ["/app/auth-service"]