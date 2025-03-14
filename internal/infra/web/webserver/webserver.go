package webserver

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//ele recebe nossas rotas usando chi
//ele recebe um mapa de handlers passo uma string de url e qual handle vou utilizar

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

// cria um novo webserver
func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

// passo o caminhoe  qual handler sera executado qdo eu apsar /order ele executara o Create do order_handler
func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

// Start inicia o servidor web pegando o middleware de logs do chi
func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	//aqui ele percorre todos os handlers que eu adicionei e adiciona no meu router
	for path, handler := range s.Handlers {
		s.Router.HandleFunc(path, handler)
	}

	http.ListenAndServe(s.WebServerPort, s.Router)

	return nil
}

func (s *WebServer) AddHealthCheck() {
	s.Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func (s *WebServer) Get(path string, handler http.HandlerFunc) {
	s.Router.MethodFunc("GET", path, handler)
}
