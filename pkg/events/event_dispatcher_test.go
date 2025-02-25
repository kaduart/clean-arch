package events

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// cria o evento
type TestEvent struct {
	Name    string      `json:"name"`
	Payload interface{} `json:"payload"`
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

/*
	func (e *TestEvent) GetDateTime() time.Time {
		return time.Now()
	}
*/
func (e *TestEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

// cria o handler
type TestEventHandler struct {
	ID int
}

// apos passar a go routine e inserir o waitGroup e mudarmpos a nossa interface temos tb que alterar os testes entao precisaremos alterar o handle com o wg
func (h *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) error {
	return nil
}

type EventDispatcherTestSuite struct {
	//usaremos o suite case para fazer testes e nao ficar com codigo repetitivo - usando testfy
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

// para as propriedades do eventDispatcherTestSuite estarem preparadas qdo eu for executar os testes eu crio o SetupTest(ele roda entre cada metodo que a gente trabalha)
func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandler{
		ID: 1,
	}
	suite.handler2 = TestEventHandler{
		ID: 2,
	}
	suite.handler2 = TestEventHandler{
		ID: 3,
	}
	suite.event = TestEvent{Name: "Test Event", Payload: "Payload"}
	suite.event2 = TestEvent{Name: "Test Event 2", Payload: "Payload 2"}
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	// handler/registro 1 - dou um registrar passo nome do evento e registro o handler
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	// para ver a qtd de eventos registrados  no meu dispatcher é 1, pois ate aqui registrei so umque foi o acima
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	// handler/registro 2 -
	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	// para ver a qtd de eventos registrados  no meu dispatcher é 2, pois ate aqui registrei dois que foram os acima
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	// prova real para saber se o hanlder que a gente registrou la e o mesmo que pedimos para ser registrado,
	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.handlers[suite.event.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][1])

	// testando registrar um evento que ja existe

}

// funcao para cecar se existe regirstro
 func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_withSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Equal(ErrHandlerAlreadyExists, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	// event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	//event 2
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	// apagando o handler do event 1
	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	// event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	//verifica se os eventos estao registrados la
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
	// o hanldr3 tem que retornar falso pois eu nao registrei ele acima somente o handler 1 e 2
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler3))
} 

type MockHandler struct {
	mock.Mock
}

// crio um metodod handle e apos inserir wg e go routine temos que passar tb o wg aqui
func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done() //obrigado passsar o done se nao ele nao vai terminar nunca a execucao esperando para sempre
}

// metodo que executa e testa todos meus eventos
func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	//trabalharemos com nosso handler e com mocks
	eh := &MockHandler{}
	eh.On("Handle", &suite.event)

	eh2 := &MockHandler{}
	eh2.On("Handle", &suite.event)
	// nao retorna nada
	/* suite.eventDispatcher.Register(suite.event.GetName(), eh)  //registrando meu mock handler passando nome evento e o mock
	suite.eventDispatcher.Register(suite.event.GetName(), eh2)  */ //registrando meu mock handler passando nome evento e o mock

	suite.eventDispatcher.Dispatch(&suite.event) //executar o dispatcher e garantir que o handle foi executado

	eh.AssertExpectations(suite.T()) // garantir que o mock foi chamado com o evento que eu disparo
	eh.AssertExpectations(suite.T()) // garantir que o mock foi chamado com o evento que eu disparo
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)
	eh.AssertNumberOfCalls(suite.T(), "Handle", 1)

}

/* func (suite *EventDispatcherTestSuite) TesteEventDispatcher_Remove() {
	// event 1 - dou um registrar passo nome do evento e registro o handler
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	// espera que ele enconntre dois registros(1)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	// espera que ele enconntre dois registros
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	// event 2
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	// handler/remover 1 - dou um remover passo nome do evento e handler
	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	//apos remover fica apenas um eVENTO -po r que tinham dois
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	// garantido se removeu o handlers correto
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.handlers[suite.event.GetName()][0])

	// handler/remover 2 -
	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler2)
	//apos remover fica apenas 0 eVENTO -po r que tinham 1
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	// handler/remover 3 -
	suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler3)
	suite.Equal(0, len(suite.eventDispatcher.handlers[suite.event2.GetName()]), err)
} */

// pedimos ao go para rodar nossa suite de test
func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
