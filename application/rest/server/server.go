package server

import (
	"log"
	"net/http"

	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	router      *mux.Router
	snipHandler *handlers.SnipHandler
}

func New(snipHandler *handlers.SnipHandler) *HttpServer {
	return &HttpServer{
		snipHandler: snipHandler,
	}
}

func (s *HttpServer) InitRouter() {
	s.router = mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./hci/static/"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	s.router.HandleFunc("/", s.snipHandler.Home).Methods("GET")
	s.router.HandleFunc("/snip/view", s.snipHandler.SnipView).Methods("GET")
	s.router.HandleFunc("/snip/create", s.snipHandler.SnipCreate).Methods("POST")
}

func (s *HttpServer) Start() {

	log.Println("Server Starting on :3000")
	err := http.ListenAndServe(":3000", s.router)
	log.Fatal(err)
}
