# ==============================================================================
# Stage 1: Build
# ==============================================================================
FROM golang:1.26-bookworm AS builder
WORKDIR /build

# Install C dependencies for cgo (ACBrLib needs OpenSSL + LibXml2)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libssl-dev \
    libxml2-dev \
    libgtk-3-dev \
    libgdk-pixbuf-2.0-dev \
    libpango1.0-dev \
    libcairo2-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Copy go module files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
COPY lib/libacbrnfe64.so /build/lib/

ENV LD_LIBRARY_PATH=/build/lib

# Build the application
RUN mkdir -p /build/bin && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /build/bin/goacbr-api ./cmd/api

# ==============================================================================
# Stage 2: Runtime
# ==============================================================================
FROM debian:bookworm-slim
WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libssl3 \
    libxml2 \
    curl \
    libgtk-3-0 \
    libgdk-pixbuf-2.0-0 \
    libpango-1.0-0 \
    libcairo2 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r goacbr && useradd -r -g goacbr -d /app -s /sbin/nologin goacbr

# Create required directories
RUN mkdir -p /app/lib /app/data/Schemas /app/logs /app/tmp && \
    chown -R goacbr:goacbr /app

# Copy ACBrLib shared library
COPY --chown=goacbr:goacbr lib/libacbrnfe64.so /app/lib/
COPY lib/libacbrnfe64.so /usr/local/lib/
RUN ldconfig

# Copy compiled binary from builder
COPY --from=builder --chown=goacbr:goacbr /build/bin/goacbr-api /app/goacbr-api

# Copy migrations (REMOVIDO: agora embutido no binário)

# Set library path for ACBrLib
ENV LD_LIBRARY_PATH=/app/lib

# Switch to non-root user
USER goacbr

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/goacbr-api"]