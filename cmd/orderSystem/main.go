package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devfull/25-clean-architeture/configs"
	"github.com/devfull/25-clean-architeture/internal/infra/graph"
	"github.com/devfull/25-clean-architeture/internal/infra/grpc/protofiles/pb"
	"github.com/devfull/25-clean-architeture/internal/infra/grpc/service"
	"github.com/devfull/25-clean-architeture/internal/infra/web"
	"github.com/devfull/25-clean-architeture/internal/infra/web/webserver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/mysql_clean_arch?parseTime=true&timeout=30s&charset=utf8mb4&collation=utf8mb4_unicode_ci")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected to database successfully")

	_, err = db.Exec(` CREATE TABLE IF NOT EXISTS orders (
            id VARCHAR(255) PRIMARY KEY,
            price DECIMAL(10,2) NOT NULL,
            tax DECIMAL(10,2) NOT NULL,
            final_price DECIMAL(10,2) NOT NULL
        ) ENGINE=InnoDB`)
	if err != nil {
		log.Fatalf("Falha ao criar tabela: %v", err)
	}

	rabbitMQChannel := getRabbitMQChannel()

	createOrderUseCase := InitializeCreateOrderUseCase(db, rabbitMQChannel)
	listOrdersUseCase := InitializeListOrdersUseCase(db)

	webserver := webserver.NewWebServer("0.0.0.0:" + configs.WebServerPort)
	fmt.Printf("Starting > web server on port: %s\n", configs.WebServerPort)
	webOrderHandler := web.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.ListOrders)

	go func() {
		if err := webserver.Start(); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Printf("Starting >> gRPC server on port:%s\n", configs.GRPCServerPort)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:"+configs.GRPCServerPort))
	if err != nil {
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("Failed to start gRPC Server: %v", err)
			}
		}()
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Printf("Starting >>> GraphQL server on port:%s\n", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp091.Channel {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to RabbitMQ successfully")
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	fmt.Println("Channel created successfully")
	return ch
}
