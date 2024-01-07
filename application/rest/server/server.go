package server

import (
	"net/http"

	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/bersennaidoo/arcbox/application/rest/mid"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/spf13/viper"
)

type HttpServer struct {
	router      *mux.Router
	snipHandler *handlers.SnipHandler
	config      *viper.Viper
	log         *golog.Logger
	mid         *mid.Middleware
}

func New(snipHandler *handlers.SnipHandler, config *viper.Viper, log *golog.Logger, mid *mid.Middleware) *HttpServer {
	return &HttpServer{
		router:      mux.NewRouter(),
		snipHandler: snipHandler,
		config:      config,
		log:         log,
		mid:         mid,
	}
}

func (s *HttpServer) InitRouter() {

	fileServer := http.FileServer(http.Dir("./hci/static/"))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	s.router.Use(s.mid.RecoverPanic, s.mid.LogRequest, s.mid.SecureHeaders)
	s.router.HandleFunc("/", s.snipHandler.Home).Methods("GET")
	s.router.HandleFunc("/snip/view/{id:[0-9]+}", s.snipHandler.SnipView).Methods("GET")
	s.router.HandleFunc("/snip/create", s.snipHandler.SnipCreate).Methods("GET")
	s.router.HandleFunc("/snip/create", s.snipHandler.SnipCreatePost).Methods("POST")
}

func (s *HttpServer) Start() {

	addr := s.config.GetString("http.http_addr")
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	s.log.Debugf("Server Starting on :3000")

	err := srv.ListenAndServe()
	if err != nil {
		s.log.Fatal(err)
	}
}
