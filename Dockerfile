# Etapa de build
FROM golang:alpine AS builder

RUN apk add --no-cache wget tar git
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize*.tar.gz \
    && rm dockerize*.tar.gz

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server -ldflags="-s -w" ./cmd/orderSystem
RUN apk add --no-cache curl gcc musl-dev

# Etapa final
FROM alpine:3.18
RUN apk add --no-cache curl
WORKDIR /app

COPY --from=builder /server /app/server
COPY cmd/orderSystem/.env /app/.env

# Copia também o dockerize
COPY --from=builder /usr/local/bin/dockerize /usr/local/bin/
RUN chmod +x /app/server
CMD ["dockerize", "-wait", "tcp://mysql:3306", "-wait", "tcp://rabbitmq:5672", "-timeout", "300s", "/app/server"]
