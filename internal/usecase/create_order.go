package usecase

import (
	"fmt"
	"log"

	"github.com/devfull/25-clean-architeture/internal/entity"
	"github.com/devfull/25-clean-architeture/pkg/events"
)

type OrderInputDTO struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
	EventDispatcher   events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	repo entity.OrderRepositoryInterface,
	dispatcher events.EventDispatcherInterface,
	event events.EventInterface,
	orderCreatedHandler events.EventHandlerInterface) *CreateOrderUseCase {

	err := dispatcher.Register(event.GetName(), orderCreatedHandler)
	if err != nil {
		panic(fmt.Sprintf("Failed to register handler: %v", err))
	}

	return &CreateOrderUseCase{
		OrderRepository:   repo,
		EventDispatcher:   dispatcher,
		OrderCreatedEvent: event,
	}
}

func (uc *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {

	order := entity.Order{
		ID:    input.ID,
		Price: input.Price,
		Tax:   input.Tax,
	}

	order.CalculateFinalPrice()
	if err := uc.OrderRepository.Save(&order); err != nil {
		log.Printf("error dispatching event: %v", err)
		return OrderOutputDTO{}, err
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	log.Println("Order event dispatched successfully")
	if uc.OrderCreatedEvent == nil {
		log.Println("OrderCreatedEvent is nil")
		return dto, fmt.Errorf("event not initialized")
	}
	uc.OrderCreatedEvent.SetPayload(dto)

	if err := uc.EventDispatcher.Dispatch(uc.OrderCreatedEvent); err != nil {
		fmt.Errorf("error dispatching event: %v", err)
		return dto, err
	}

	return dto, nil
}
