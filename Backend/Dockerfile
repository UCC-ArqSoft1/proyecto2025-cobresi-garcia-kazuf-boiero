# Build stage: compile the Go backend.
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o backend ./cmd/server

# Runtime stage: minimal image to run the compiled binary.
FROM alpine:3.20 AS runtime
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H -s /sbin/nologin appuser

COPY --from=builder /app/backend ./backend

EXPOSE 8080
USER appuser
CMD ["./backend"]
