# Dockerfile
# Build stage
FROM golang:1.24.5-alpine AS builder
WORKDIR /app

# Install git for Go modules
RUN apk add --no-cache git

# Enable static build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Cache go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source code
COPY . .

# Build binary
RUN --mount=type=cache,target=/root/.cache/go-build go build -o server .

# Final stage
FROM alpine:3.20
WORKDIR /app

# Non-root user
RUN adduser -D appuser
USER appuser

# Copy binary from builder
COPY --from=builder /app/server .

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]


