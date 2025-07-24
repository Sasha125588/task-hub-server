FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Install required system packages
RUN apk add --no-cache gcc musl-dev git

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/app/main.go

# Build the application
RUN go build -o app ./cmd/app

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required runtime packages
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"] 