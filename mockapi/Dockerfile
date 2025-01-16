# Etapa 1: Construção da aplicação mockapi
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mockapi

# Etapa 2: Criar a imagem final para mockapi
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mockapi .

EXPOSE 8000

CMD ["./mockapi"]