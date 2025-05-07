FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o life_log ./cmd/life-log

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/life_log /app/life_log
EXPOSE 8080
ENTRYPOINT ["/app/life_log"]
