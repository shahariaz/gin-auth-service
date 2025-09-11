# Multi-stage build for efficiency
FROM golang:1.23-alpine AS builder  
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app (disable cgo for smaller binary, static linking)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Final stage: Runtime image (Alpine for small size)
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env and static files (if any)
COPY .env .
COPY static ./static  
# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the app
CMD ["./main"]