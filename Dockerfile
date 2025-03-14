FROM golang:alpine as builder

RUN apk add --no-cache wget tar git
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize*.tar.gz \
    && rm dockerize*.tar.gz

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server  # Otimização de binário
RUN apk add --no-cache curl gcc musl-dev

FROM rabbitmq:3.13.7-management

COPY rabbitmq.conf /etc/rabbitmq/
COPY definitions.json /etc/rabbitmq/
COPY init.sh /docker-entrypoint-init.d/

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh
RUN chmod +x /docker-entrypoint-init.d/init.sh

FROM alpine:3.18
RUN apk add --no-cache curl
COPY --from=builder /server /server
COPY --from=builder /usr/local/bin/dockerize /usr/local/bin/

CMD ["dockerize", "-wait", "tcp://mysql:3306", "-wait", "tcp://rabbitmq:5672", "-timeout", "300s", "/server"]
