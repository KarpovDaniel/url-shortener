FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o url-shortener ./cmd/url-shortener

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/url-shortener .
COPY .env .
COPY migrations/ ./migrations/
EXPOSE 8080 50051
CMD ["./url-shortener"]