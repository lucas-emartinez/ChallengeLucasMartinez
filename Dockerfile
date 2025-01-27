# Build stage
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod tidy

# Copy all source code and start.sh
COPY . .
COPY start.sh ./start.sh

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker ./cmd/worker/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/consumer ./cmd/consumer/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binaries and start.sh
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/bin/worker /app/worker
COPY --from=builder /app/bin/consumer /app/consumer
COPY --from=builder /app/start.sh /app/start.sh

# Make sure start.sh is executable
RUN chmod +x /app/start.sh

EXPOSE 8080

# Use full path to start.sh
CMD ["/app/start.sh"]