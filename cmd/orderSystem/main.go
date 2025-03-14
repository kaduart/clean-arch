package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=30s&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		configs.DBUser,
		configs.DBPassword,
		configs.DBHost,
		configs.DBPort,
		configs.DBName)
	db, err := sql.Open(configs.DBDriver, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Connected to database successfully")

	rabbitMQChannel := getRabbitMQChannel()

	createOrderUseCase := InitializeCreateOrderUseCase(db, rabbitMQChannel)
	listOrdersUseCase := InitializeListOrdersUseCase(db)

	webserver := webserver.NewWebServer("0.0.0.0:" + configs.WebServerPort)
	fmt.Printf("Starting > web server on port: %s\n", configs.WebServerPort)

	webserver.AddHealthCheck()
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

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", configs.GRPCServerPort))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", configs.GRPCServerPort, err)
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC Server: %v", err)
		}
	}()

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	go func() {
		fmt.Printf("Starting >>> GraphQL server on port:%s\n", configs.GraphQLServerPort)
		if err := http.ListenAndServe(":"+configs.GraphQLServerPort, nil); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}

func getRabbitMQChannel() *amqp091.Channel {
	var conn *amqp091.Connection
	var err error
	for i := 0; i < 10; i++ { // Tenta conectar 10 vezes
		conn, err = amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		}
		fmt.Printf("Tentando conectar ao RabbitMQ... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second) // Espera 5 segundos antes de tentar novamente
	}
	if err != nil {
		panic("Não foi possível conectar ao RabbitMQ: " + err.Error())
	}
	fmt.Println("Connected to RabbitMQ successfully")

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	fmt.Println("Channel created successfully")

	return ch
}
