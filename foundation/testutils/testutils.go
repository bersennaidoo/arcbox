package testutils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bersennaidoo/arcbox/application/rest/application"
	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/bersennaidoo/arcbox/application/rest/mid"
	"github.com/bersennaidoo/arcbox/foundation/formdecode"
	"github.com/bersennaidoo/arcbox/infrastructure/repositories/mysql"
	"github.com/bersennaidoo/arcbox/physical/config"
	"github.com/bersennaidoo/arcbox/physical/dbc"
	"github.com/bersennaidoo/arcbox/physical/logger"
	"github.com/bersennaidoo/arcbox/physical/session"
	"github.com/gorilla/mux"
)

func NewTestApplication(t *testing.T) *application.Application {
	log := logger.New()
	filename := config.GetConfigFileName()
	cfg := config.New(filename)
	dbc := dbc.New(cfg)
	sr := mysql.NewSnipsRepository(dbc)
	ur := mysql.NewUsersRepository(dbc)
	decoder := formdecode.New()
	tempcache, _ := handlers.NewTemplateCache()
	sessionM := session.New(dbc)
	h := handlers.New(log, sr, ur, tempcache, decoder, sessionM)
	m := mid.New(log, sessionM, ur)

	return &application.Application{
		Router:         mux.NewRouter(),
		Handlers:       h,
		Config:         cfg,
		Log:            log,
		Mid:            m,
		SessionManager: sessionM,
	}
}

type TestServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewTLSServer(h)
	return &TestServer{ts}
}

func (ts *TestServer) Get(t *testing.T, urlPath string) (int, http.Header, string) {
	t.Helper()
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
