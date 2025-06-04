# Build stage
FROM golang:1.24.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates build-base

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .


# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s -extldflags "-static"' -o /app/server ./cmd

# Final stage
FROM alpine:3.20

# Install runtime dependencies
RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app

# Copy binary and configs
COPY --from=builder /app/server .
COPY --from=builder /app/configs/ ./configs/
COPY --from=builder /app/i18n/ ./i18n/
COPY --from=builder /app/casbin/ ./casbin/

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

# Expose port
EXPOSE 3033

# Run the application
CMD ["./server"]
