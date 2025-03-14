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
	"github.com/devfull/25-clean-architeture/internal/usecase"
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

	startWebServer(configs, createOrderUseCase, listOrdersUseCase)

	startGRPC(createOrderUseCase, listOrdersUseCase, configs)

	startGraphQL(createOrderUseCase, listOrdersUseCase, configs)

	select {}
}

func startWebServer(cfg *configs.Conf, createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) {
	webserver := webserver.NewWebServer("0.0.0.0:" + cfg.WebServerPort)
	fmt.Printf("Starting > web server on port: %s\n", cfg.WebServerPort)

	webserver.AddHealthCheck()
	webOrderHandler := web.NewWebOrderHandler(createOrderUseCase, listOrdersUseCase)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.ListOrders)

	go func() {
		if err := webserver.Start(); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	}()
}

func startGRPC(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase, cfg *configs.Conf) {
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", cfg.GRPCServerPort))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", cfg.GRPCServerPort, err)
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC Server: %v", err)
		}
	}()
}

func startGraphQL(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase, cfg *configs.Conf) {
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	go func() {
		fmt.Printf("Starting >>> GraphQL server on port:%s\n", cfg.GraphQLServerPort)
		if err := http.ListenAndServe(":"+cfg.GraphQLServerPort, nil); err != nil {
			log.Fatal(err)
		}
	}()
}

func getRabbitMQChannel() *amqp091.Channel {
	var conn *amqp091.Connection
	var err error
	for i := 0; i < 10; i++ {
		conn, err = amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		}
		fmt.Printf("Tentando conectar ao RabbitMQ... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		panic("Não foi possível conectar ao RabbitMQ: " + err.Error())
	}
	fmt.Println("Connected to RabbitMQ successfully")

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	ch.Qos(10, 0, false)
	fmt.Println("Channel created successfully")

	return ch
}
