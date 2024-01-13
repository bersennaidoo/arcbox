package server

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/bersennaidoo/arcbox/application/rest/mid"
	"github.com/bersennaidoo/arcbox/hci"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/spf13/viper"
)

type HttpServer struct {
	router         *mux.Router
	snipHandler    *handlers.SnipHandler
	config         *viper.Viper
	log            *golog.Logger
	mid            *mid.Middleware
	sessionManager *scs.SessionManager
}

func New(snipHandler *handlers.SnipHandler, config *viper.Viper, log *golog.Logger,
	mid *mid.Middleware, sessionManager *scs.SessionManager) *HttpServer {
	return &HttpServer{
		router:         mux.NewRouter(),
		snipHandler:    snipHandler,
		config:         config,
		log:            log,
		mid:            mid,
		sessionManager: sessionManager,
	}
}

func (s *HttpServer) InitRouter() {

	//fileServer := http.FileServer(http.Dir("./hci/static/"))
	//s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	fileServer := http.FileServer(http.FS(hci.Files))
	s.router.PathPrefix("/static").Handler(http.StripPrefix("", fileServer))

	auth := s.router.PathPrefix("/snip").Subrouter()
	authu := s.router.PathPrefix("/user").Subrouter()

	s.router.Use(s.mid.RecoverPanic, s.mid.LogRequest, s.mid.Authenticate, s.mid.SecureHeaders)
	auth.Use(s.mid.RecoverPanic, s.mid.LogRequest, s.mid.RequireAuthentication, s.mid.SecureHeaders)
	authu.Use(s.mid.RecoverPanic, s.mid.LogRequest, s.mid.RequireAuthentication, s.mid.SecureHeaders)

	auth.HandleFunc("/create", s.snipHandler.SnipCreate).Methods("GET")
	auth.HandleFunc("/create", s.snipHandler.SnipCreatePost).Methods("POST")

	authu.HandleFunc("/logout", s.snipHandler.UserLogoutPost).Methods("POST")

	s.router.HandleFunc("/", s.snipHandler.Home).Methods("GET")
	s.router.HandleFunc("/snip/view/{id:[0-9]+}", s.snipHandler.SnipView).Methods("GET")
	s.router.HandleFunc("/user/signup", s.snipHandler.UserSignup).Methods("GET")
	s.router.HandleFunc("/user/signup", s.snipHandler.UserSignupPost).Methods("POST")
	s.router.HandleFunc("/user/login", s.snipHandler.UserLogin).Methods("GET")
	s.router.HandleFunc("/user/login", s.snipHandler.UserLoginPost).Methods("POST")
	http.Handle("/", s.router)
}

func (s *HttpServer) Start() {

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}

	addr := s.config.GetString("http.http_addr")
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.sessionManager.LoadAndSave(s.router),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.log.Debugf("Server Starting on :3000")

	err := srv.ListenAndServeTLS("./documentation/certs/cert.pem",
		"./documentation/certs/key.pem")
	if err != nil {
		s.log.Fatal(err)
	}
}
