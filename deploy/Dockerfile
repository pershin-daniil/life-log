FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o life_log ./cmd/life-log

FROM alpine:latest
COPY --from=builder /app/life_log /app
EXPOSE 8080
ENTRYPOINT ["/app"]
