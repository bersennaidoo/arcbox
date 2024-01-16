package application

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/bersennaidoo/arcbox/application/rest/mid"
	"github.com/bersennaidoo/arcbox/hci"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/kataras/golog"
	"github.com/spf13/viper"
)

type Application struct {
	Handlers       *handlers.Handler
	Config         *viper.Viper
	Log            *golog.Logger
	Mid            *mid.Middleware
	SessionManager *scs.SessionManager
}

func New(Handlers *handlers.Handler, Config *viper.Viper, Log *golog.Logger,
	Mid *mid.Middleware, SessionManager *scs.SessionManager) *Application {
	return &Application{
		Handlers:       Handlers,
		Config:         Config,
		Log:            Log,
		Mid:            Mid,
		SessionManager: SessionManager,
	}
}

func (a *Application) InitRouter() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Handlers.NotFound(w)
	})

	fileServer := http.FileServer(http.FS(hci.Files))

	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", handlers.Ping)

	dynamic := alice.New(a.SessionManager.LoadAndSave, a.Mid.Authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(a.Handlers.Home))
	router.Handler(http.MethodGet, "/snip/view/:id", dynamic.ThenFunc(a.Handlers.SnipView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(a.Handlers.UserSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(a.Handlers.UserSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(a.Handlers.UserLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(a.Handlers.UserLoginPost))

	protected := dynamic.Append(a.Mid.RequireAuthentication)

	router.Handler(http.MethodGet, "/snip/create", protected.ThenFunc(a.Handlers.SnipCreate))
	router.Handler(http.MethodPost, "/snip/create", protected.ThenFunc(a.Handlers.SnipCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(a.Handlers.UserLogoutPost))

	standard := alice.New(a.Mid.RecoverPanic, a.Mid.LogRequest, a.Mid.SecureHeaders)
	return standard.Then(router)
}

func (a *Application) Start(n http.Handler) {

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}

	addr := a.Config.GetString("http.http_addr")
	srv := &http.Server{
		Addr:         addr,
		Handler:      a.SessionManager.LoadAndSave(n),
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
