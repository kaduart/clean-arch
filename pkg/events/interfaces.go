package events

import (
	"sync"
)

/*
	estrutura:
		evento(carrega dados)
		-> Operacoes que serao executadas qdo um evento e chamado
		-> Gerenciador dos nossos eventos/operacoes: responsavel por
			-> Registrar os eventos e suas operacoes
			-> Despachar/Fire no evento para suas operacoes sejam executadas

*/

// evento
type EventInterface interface {
	GetName() string
	GetPayload() interface{}
	SetPayload(payload interface{})
}

// operacao que executa o evento
type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup) error //metodo para tratar/executar nossa operacao a partir de um evento
}

// Registrar os eventos e suas operacoes
type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error // quando o eventName for executado executa o handler
	Dispatch(event EventInterface) error                            // despacha(dar um fire) no evento
	Remove(eventName string, handler EventHandlerInterface) error   //
	Has(eventName string, handler EventHandlerInterface) bool       //
	Clear() error                                                   // limpa nosso dispatcher matando todos eventos que tem
}
