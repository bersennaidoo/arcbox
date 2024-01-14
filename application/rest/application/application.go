package application

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

type Application struct {
	router         *mux.Router
	Handlers       *handlers.Handler
	config         *viper.Viper
	log            *golog.Logger
	mid            *mid.Middleware
	sessionManager *scs.SessionManager
}

func New(Handlers *handlers.Handler, config *viper.Viper, log *golog.Logger,
	mid *mid.Middleware, sessionManager *scs.SessionManager) *Application {
	return &Application{
		router:         mux.NewRouter(),
		Handlers:       Handlers,
		config:         config,
		log:            log,
		mid:            mid,
		sessionManager: sessionManager,
	}
}

func (a *Application) InitRouter() {

	fileServer := http.FileServer(http.FS(hci.Files))
	a.router.PathPrefix("/static").Handler(http.StripPrefix("", fileServer))

	auth := a.router.PathPrefix("/snip").Subrouter()
	authu := a.router.PathPrefix("/user").Subrouter()

	a.router.Use(a.mid.RecoverPanic, a.mid.LogRequest, a.mid.Authenticate, a.mid.SecureHeaders)
	auth.Use(a.mid.RecoverPanic, a.mid.LogRequest, a.mid.RequireAuthentication, a.mid.SecureHeaders)
	authu.Use(a.mid.RecoverPanic, a.mid.LogRequest, a.mid.RequireAuthentication, a.mid.SecureHeaders)

	auth.HandleFunc("/create", a.Handlers.SnipCreate).Methods("GET")
	auth.HandleFunc("/create", a.Handlers.SnipCreatePost).Methods("POST")

	authu.HandleFunc("/logout", a.Handlers.UserLogoutPost).Methods("POST")

	a.router.HandleFunc("/", a.Handlers.Home).Methods("GET")
	a.router.HandleFunc("/snip/view/{id:[0-9]+}", a.Handlers.SnipView).Methods("GET")
	a.router.HandleFunc("/user/signup", a.Handlers.UserSignup).Methods("GET")
	a.router.HandleFunc("/user/signup", a.Handlers.UserSignupPost).Methods("POST")
	a.router.HandleFunc("/user/login", a.Handlers.UserLogin).Methods("GET")
	a.router.HandleFunc("/user/login", a.Handlers.UserLoginPost).Methods("POST")
	http.Handle("/", a.router)
}

func (a *Application) Start() {

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}

	addr := a.config.GetString("http.http_addr")
	srv := &http.Server{
		Addr:         addr,
		Handler:      a.sessionManager.LoadAndSave(a.router),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a.log.Debugf("Server Starting on :3000")

	err := srv.ListenAndServeTLS("./documentation/certs/cert.pem",
		"./documentation/certs/key.pem")
	if err != nil {
		a.log.Fatal(err)
	}
}
