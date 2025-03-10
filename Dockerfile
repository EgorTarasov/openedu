FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code, excluding extension and bruno folders
COPY . .
# Remove unwanted directories 
RUN rm -rf extension/ bruno/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main cmd/server/main.go

# Runner stage
FROM alpine:3.21

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy views directory for HTML templates
COPY --from=builder /app/views ./views

# Expose application port (adjust as needed)
EXPOSE 8080

# Run the application
CMD ["/app/main"]