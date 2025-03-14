Clean Architecture Go Application
A multi-protocol API for order management, implemented using Clean Architecture principles.
The application offers endpoints via gRPC, HTTP, and GraphQL, and also uses RabbitMQ for event handling and MySQL for persistence.

Project Structure
plaintext
Copy
Edit
.
├── cmd/
│   └── orderSystem/
│       └── main.go          # Application entry point
├── internal/
│   ├── infra/
│   │   ├── database/        # Repositories and MySQL connections
│   │   ├── grpc/            # gRPC configuration and implementation
│   │   └── http/            # HTTP handlers and server configuration (REST)
├── migrations/              # Database migration scripts
│   └── 0001_create_orders.sql
├── scripts/                 # Auxiliary scripts
│   └── wait-for-db.sh       # Waits for the database to be ready
├── docker-compose.yml       # Docker Compose configuration (MySQL, RabbitMQ, etc.)
├── Dockerfile               # Dockerfile for building the application
├── go.mod                   # Go dependencies management
└── README.md                # Documentation (this file)
Features
Order Management
Create orders with automatic final price calculation.
List all registered orders.
Multi-Protocol Communication
gRPC: Available on port 50051.
REST HTTP: Available on port 8000.
GraphQL: Available on port 8083.
Event Handling
Integration with RabbitMQ for event dispatching and processing.
Persistence
MySQL database for storing order data.
Technologies Used
Language: Go 1.20+
gRPC: High-performance RPC protocol
GraphQL: Implementation using gqlgen
HTTP Router: Chi Router
Dependency Injection: Wire
Database: MySQL 8+
Message Broker: RabbitMQ
Configuration Management: Viper
Installation
Prerequisites
Go 1.20+
Docker and Docker Compose
MySQL 8+ (if not using the container)
RabbitMQ (if not using the container)
Protoc (for compiling protobuf files)
Steps
Clone the repository:

bash
Copy
Edit
git clone https://github.com/your-username/clean-arch-go.git
cd clean-arch-go
Docker Compose Configuration

The docker-compose.yml file is pre-configured to start the required containers:

yaml
Copy
Edit
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
    environment:
      # Example custom variable if needed:
      # RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS: "-rabbit handshake_timeout 60000"
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

volumes:
  mysql_data:
  rabbitmq_data:
Environment Configuration

Create a .env file in the project root with the following variables:

dotenv
Copy
Edit
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=clean_arch

# Ports
GRPC_PORT=50051
HTTP_PORT=8000
GRAPHQL_PORT=8083

# RabbitMQ
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
.gitignore File

Ensure that unnecessary files and directories are ignored by Git. For example:

gitignore
Copy
Edit
# Binaries
bin/
vendor/
*.exe
*.dll
*.so
*.dylib

# Dependencies
go.sum

# Environment
.env
*.env.local

# Docker
.docker/mysql/
docker-compose.override.yml

# Tests
*.test
*.prof
coverage.txt

# IDEs
.idea/
.vscode/
*.swp
*.swo

# Logs and temporary files
*.log
*.tmp
Dependency Injection with Wire

The Wire configuration is in cmd/orderSystem/wire.go and defines how components are injected:

go
Copy
Edit
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
Project Dependencies (go.mod)

go
Copy
Edit
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
GraphQL Configuration (gqlgen.yml)

yaml
Copy
Edit
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
Execution
Using Docker Compose
Start the Containers:

bash
Copy
Edit
docker-compose up -d
Wait for MySQL and RabbitMQ Initialization:
You can use the script scripts/wait-for-db.sh to wait until the database is ready before starting the application.

Run the Application Locally (for development):

bash
Copy
Edit
go run ./cmd/orderSystem/main.go
Testing Endpoints
gRPC:
Use Evans or any other gRPC tool to test the service on port 50051.

HTTP (REST):
Make requests to http://localhost:8000/health to check the API’s health.

GraphQL:
Access http://localhost:8083 to use the GraphQL interface and test queries/mutations.

Final Considerations
Migrations:
Use the scripts in the migrations/ folder to create or update the database schema.

Logs and Debugging:
Monitor the container logs to check the status and behavior of the services (RabbitMQ, MySQL, and the application).

Customization:
Adjust the environment variables in the .env file according to your development or production environment needs.