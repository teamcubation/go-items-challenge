FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mockapi

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mockapi .

EXPOSE 8000

CMD ["./mockapi"]