package testutils

import (
	"bytes"
	"html"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bersennaidoo/arcbox/foundation/formdecode"
	"github.com/bersennaidoo/arcbox/infrastructure/repositories/mocks"
	"github.com/bersennaidoo/arcbox/physical/config"
	"github.com/bersennaidoo/arcbox/physical/logger"
	"github.com/bersennaidoo/arcbox/transport/http/application"
	"github.com/bersennaidoo/arcbox/transport/http/handlers"
	"github.com/bersennaidoo/arcbox/transport/http/mid"
)

func NewTestApplication(t *testing.T) *application.Application {
	log := logger.New()
	filename := config.GetConfigFileName()
	cfg := config.New(filename)
	//dbc := dbc.New(cfg)
	//sr := mysql.NewSnipsRepository(dbc)
	//ur := mysql.NewUsersRepository(dbc)
	msr := &mocks.SnipsMockRepository{}
	mur := &mocks.UsersMockRepository{}
	decoder := formdecode.New()
	tempcache, _ := handlers.NewTemplateCache()
	//sessionM := session.New(dbc)
	sessionMockManager := scs.New()
	sessionMockManager.Lifetime = 12 * time.Hour
	sessionMockManager.Cookie.Secure = true

	h := handlers.New(log, msr, mur, tempcache, decoder, sessionMockManager)
	//m := mid.New(log, sessionM, ur)
	midMock := mid.New(log, sessionMockManager, mur)

	return &application.Application{
		Handlers:       h,
		Config:         cfg,
		Log:            log,
		Mid:            midMock,
		SessionManager: sessionMockManager,
	}
}

type TestServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	t.Helper()
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

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

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func ExtractCSRFToken(t *testing.T, body string) string {

	matches := csrfTokenRX.FindStringSubmatch(body)

	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}
