FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/api ./cmd/api

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary
COPY --from=builder /app/bin/api .

# Copy configuration
COPY configs /app/configs

# Create a non-root user
RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./api"] 