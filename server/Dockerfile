# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build binary named 'server'
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Stage 2: Runtime (Tiny image)
FROM alpine:latest
WORKDIR /root/
# Install minimal dependencies (CA certs for HTTPS calls)
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]