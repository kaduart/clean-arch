Structure:
25.1-CLEAN-ARCHITETURE-APP/
├── cmd/
│   └── orderSystem/
│       └── wire.go
├── internal/
│   ├── infra/
│   │   ├── database/
│   │   ├── grpc/
│   │   └── graphql/
├── protofiles/
├── .env
├── .gitignore
├── docker-compose.yml
├── go.mod
└── Makefile

# Clean Architecture Go Application

API multi-protocol (gRPC, HTTP, GraphQL) para gestão de pedidos seguindo princípios de Clean Architecture

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## 📋 Funcionalidades

- Criação de pedidos com cálculo automático de preço final
- Listagem de todos os pedidos cadastrados
- Comunicação via múltiplos protocolos:
  - **gRPC** (porta 50051)
  - **REST HTTP** (porta 8000)
  - **GraphQL** (porta 8083)
- Sistema de eventos com RabbitMQ
- Banco de dados MySQL para persistência

## 🛠 Tecnologias

- **Linguagem**: Go 1.20+
- **gRPC**: Protocolo RPC de alto desempenho
- **GraphQL**: Implementação com gqlgen
- **HTTP Router**: Chi Router
- **DI**: Wire (Google)
- **Banco de Dados**: MySQL
- **Message Broker**: RabbitMQ
- **Configuração**: Viper

## 🚀 Instalação

### Pré-requisitos
- Go 1.20+
- MySQL 8+
- RabbitMQ
- Protoc (para compilação de protobuf)

### Passos
1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/clean-arch-go.git
cd clean-arch-go

1. docker-compose.yml

version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      interval: 10s
      retries: 10

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

volumes:
  mysql_data:
  rabbitmq_data:


2. .env

# Banco de Dados
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=clean_arch

# Portas
GRPC_PORT=50051
HTTP_PORT=8000
GRAPHQL_PORT=8083

# RabbitMQ
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

3. .gitignore

# Binários
bin/
vendor/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Dependências
go.sum

# Ambiente
.env
*.env.local

# Docker
.docker/mysql/
docker-compose.override.yml

# Testes
*.test
*.prof
coverage.txt

# IDEs
.idea/
.vscode/
*.swp
*.swo

# Logs e arquivos temporários
*.log
*.tmp


4. wire.go (DI Configuration)

// cmd/orderSystem/wire.go

//go:build wireinject

package main

import (
    "database/sql"
    "github.com/devfull/25-clean-architeture/internal/infra/database"
    "github.com/devfull/25-clean-architeture/internal/infra/event"
    "github.com/devfull/25-clean-architeture/internal/infra/event/handler"
    "github.com/devfull/25-clean-architeture/internal/usecase"
    "github.com/devfull/25-clean-architeture/pkg/events"
    "github.com/google/wire"
    amqp "github.com/rabbitmq/amqp091-go"
)

var SuperSet = wire.NewSet(
    database.NewOrderRepository,
    wire.Bind(new(usecase.OrderRepositoryInterface), new(*database.OrderRepository)),
    events.NewEventDispatcher,
    event.NewOrderCreatedEvent,
    handler.NewOrderCreatedHandler,
    usecase.NewCreateOrderUseCase,
    usecase.NewListOrdersUseCase,
)

func InitializeCreateOrderUseCase(db *sql.DB, ch *amqp.Channel) *usecase.CreateOrderUseCase {
    wire.Build(SuperSet)
    return &usecase.CreateOrderUseCase{}
}

5. go.mod (Exemplo)

module github.com/devfull/25-clean-architeture

go 1.20

require (
    github.com/99designs/gqlgen v0.17.31
    github.com/google/wire v0.5.0
    github.com/rabbitmq/amqp091-go v1.8.0
    github.com/stretchr/testify v1.8.4
    google.golang.org/grpc v1.56.1
    google.golang.org/protobuf v1.31.0
)

6. gqlgen.yml (GraphQL Config)

schema:
  - graph/schema.graphqls
exec:
  filename: graph/generated.go
  package: graph
model:
  filename: graph/model/models_gen.go
  package: model
resolver:
  layout: follow-schema
  dir: graph
