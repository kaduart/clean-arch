package event

type OrderCreatedEvent struct {
	payload interface{}
}

func NewOrderCreatedEvent() *OrderCreatedEvent {
	return &OrderCreatedEvent{}
}
func (e *OrderCreatedEvent) GetName() string {
	return "OrderCreated"
}

func (e *OrderCreatedEvent) GetPayload() interface{} {
	return e.payload
}

func (e *OrderCreatedEvent) SetPayload(payload interface{}) {
	e.payload = payload
}
