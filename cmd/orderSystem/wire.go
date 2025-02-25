//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/devfull/25-clean-architeture/internal/entity"
	"github.com/devfull/25-clean-architeture/internal/infra/database"
	"github.com/devfull/25-clean-architeture/internal/infra/event"
	"github.com/devfull/25-clean-architeture/internal/infra/event/handler"
	"github.com/devfull/25-clean-architeture/internal/infra/web"
	"github.com/devfull/25-clean-architeture/internal/usecase"
	"github.com/devfull/25-clean-architeture/pkg/events"
	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setListOrdersDependency = wire.NewSet(
	setOrderRepositoryDependency,
	usecase.NewListOrdersUseCase,
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreatedEvent, // Provider do evento
	wire.Bind(new(events.EventInterface), new(*event.OrderCreatedEvent)), // Liga interface à implementação
)

var setOrderCreatedHandler = wire.NewSet(
	handler.NewOrderCreatedHandler,
	wire.Bind(new(events.EventHandlerInterface), new(*handler.OrderCreatedHandler)),
)

func InitializeCreateOrderUseCase(db *sql.DB, rabbitMQChannel *amqp091.Channel) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setEventDispatcherDependency,
		setOrderCreatedEvent,
		setOrderCreatedHandler,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func InitializeListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		setListOrdersDependency,
	)
	return &usecase.ListOrdersUseCase{}
}

func InitializeWebOrderHandler(db *sql.DB, rabbitMQChannel *amqp091.Channel) *web.WebOrderHandler {
	wire.Build(
		InitializeCreateOrderUseCase,
		InitializeListOrdersUseCase,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
