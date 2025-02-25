package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/devfull/25-clean-architeture/pkg/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

var _ events.EventHandlerInterface = (*OrderCreatedHandler)(nil)

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) error {
	defer wg.Done()

	jsonOutput, _ := json.Marshal(event.GetPayload())

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	_, err := h.RabbitMQChannel.QueueDeclare(
		"orders",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = h.RabbitMQChannel.QueueBind(
		"orders",
		"order_created",
		"amq.direct",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	err = h.RabbitMQChannel.Publish(
		"amq.direct",
		"order_created",
		false,
		false,
		msg,
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	fmt.Println("Publishing in RabbitMQ!")
	return nil
}
