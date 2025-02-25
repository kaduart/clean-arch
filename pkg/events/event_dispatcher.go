package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyExists = errors.New("Handler already exists")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
	mu       sync.RWMutex
}

var _ EventDispatcherInterface = (*EventDispatcher)(nil)

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyExists
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return errors.New("handler not found")
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	if handlers, ok := ed.handlers[event.GetName()]; ok {

		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

//criar um file para testar event_dispatcher_test

func (ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)

	return nil
}
