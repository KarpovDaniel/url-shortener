FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener

FROM alpine:latest
COPY --from=builder /app/url-shortener .
COPY .env .
COPY migrations/ ./migrations/
EXPOSE 8080
CMD ["./url-shortener"]