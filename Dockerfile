FROM golang:1.23.0-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o tender-service cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tender-service .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

CMD ["./tender-service"]