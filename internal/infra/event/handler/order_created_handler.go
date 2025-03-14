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
		return fmt.Errorf("falha ao declarar a fila 'orders': %v", err)
	}

	err = h.RabbitMQChannel.QueueBind(
		"orders",
		"order_created",
		"amq.direct",
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("falha ao vincular a fila 'orders' à troca 'amq.direct': %v", err)
	}

	err = h.RabbitMQChannel.Publish(
		"amq.direct",
		"order_created",
		false,
		false,
		msg,
	)

	if err != nil {
		return fmt.Errorf("falha ao publicar a mensagem na fila: %v", err)
	}

	fmt.Println("Mensagem publicada na fila 'orders' do RabbitMQ com sucesso!")
	return nil
}
