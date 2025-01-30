FROM golang:1.23.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copia todo el proyecto, no solo el directorio cmd/api
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM gcr.io/distroless/base-debian11
WORKDIR /
COPY --from=builder /app/api .
ENTRYPOINT ["/api"]