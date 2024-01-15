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
	Router         *mux.Router
	Handlers       *handlers.Handler
	Config         *viper.Viper
	Log            *golog.Logger
	Mid            *mid.Middleware
	SessionManager *scs.SessionManager
}

func New(Handlers *handlers.Handler, Config *viper.Viper, Log *golog.Logger,
	Mid *mid.Middleware, SessionManager *scs.SessionManager) *Application {
	return &Application{
		Router:         mux.NewRouter(),
		Handlers:       Handlers,
		Config:         Config,
		Log:            Log,
		Mid:            Mid,
		SessionManager: SessionManager,
	}
}

func (a *Application) InitRouter() (*mux.Router, *mux.Router, *mux.Router, *mux.Router) {

	fileServer := http.FileServer(http.FS(hci.Files))
	a.Router.PathPrefix("/static").Handler(http.StripPrefix("", fileServer))

	auth := a.Router.PathPrefix("/snip").Subrouter()
	authu := a.Router.PathPrefix("/user").Subrouter()

	pingr := a.Router.PathPrefix("").Subrouter()
	pingr.HandleFunc("/ping", handlers.Ping).Methods("GET")

	a.Router.Use(a.Mid.RecoverPanic, a.Mid.LogRequest, a.Mid.Authenticate, a.Mid.SecureHeaders)
	auth.Use(a.Mid.RecoverPanic, a.Mid.LogRequest, a.Mid.RequireAuthentication, a.Mid.SecureHeaders)
	authu.Use(a.Mid.RecoverPanic, a.Mid.LogRequest, a.Mid.RequireAuthentication, a.Mid.SecureHeaders)

	auth.HandleFunc("/create", a.Handlers.SnipCreate).Methods("GET")
	auth.HandleFunc("/create", a.Handlers.SnipCreatePost).Methods("POST")

	authu.HandleFunc("/logout", a.Handlers.UserLogoutPost).Methods("POST")

	a.Router.HandleFunc("/", a.Handlers.Home).Methods("GET")
	a.Router.HandleFunc("/snip/view/{id:[0-9]+}", a.Handlers.SnipView).Methods("GET")
	a.Router.HandleFunc("/user/signup", a.Handlers.UserSignup).Methods("GET")
	a.Router.HandleFunc("/user/signup", a.Handlers.UserSignupPost).Methods("POST")
	a.Router.HandleFunc("/user/login", a.Handlers.UserLogin).Methods("GET")
	a.Router.HandleFunc("/user/login", a.Handlers.UserLoginPost).Methods("POST")
	http.Handle("/", a.Router)

	return pingr, a.Router, auth, authu
}

func (a *Application) Start() {

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}

	addr := a.Config.GetString("http.http_addr")
	srv := &http.Server{
		Addr:         addr,
		Handler:      a.SessionManager.LoadAndSave(a.Router),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a.Log.Debugf("Server Starting on :3000")

	err := srv.ListenAndServeTLS("./documentation/certs/cert.pem",
		"./documentation/certs/key.pem")
	if err != nil {
		a.Log.Fatal(err)
	}
}
