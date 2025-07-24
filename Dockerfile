FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Install required system packages
RUN apk add --no-cache gcc musl-dev git

# Copy go mod files and vendor directory
COPY go.mod go.sum ./
COPY vendor/ vendor/

# Copy source code
COPY . .

# Build the application
RUN go build -mod=vendor -o /go/bin/app ./cmd/app

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required runtime packages
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /go/bin/app .
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"] 