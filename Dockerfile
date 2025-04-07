FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/url-shortener .
COPY .env .
COPY migrations/ ./migrations/
EXPOSE 8080
EXPOSE 50051
CMD ["./url-shortener"]