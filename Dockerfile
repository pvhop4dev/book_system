# Build stage
FROM golang:1.24.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o book-system ./cmd

# Final stage
FROM alpine:3.20

# Set logging driver
ENV LOG_DRIVER=local

# Install runtime dependencies
RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app

# Copy binary and necessary files only
COPY --from=builder /app/book-system .
COPY --from=builder /app/configs/ ./configs/
COPY --from=builder /app/i18n/ ./i18n/
COPY --from=builder /app/casbin/ ./casbin/
COPY --from=builder /app/config.yaml ./config.yaml

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

# Run the application with more detailed logging
CMD echo "Starting Book System..." && ./book-system