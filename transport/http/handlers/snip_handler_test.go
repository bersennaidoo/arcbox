package handlers_test

import (
	"net/http"
	"testing"

	"github.com/bersennaidoo/arcbox/foundation/assert"
	"github.com/bersennaidoo/arcbox/foundation/testutils"
)

func TestPing(t *testing.T) {

	app := testutils.NewTestApplication(t)

	arouter := app.InitRouter()

	ts := testutils.NewTestServer(t, arouter)
	defer ts.Close()

	code, _, body := ts.Get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)

	assert.Equal(t, body, "OK")
}

func TestSnipView(t *testing.T) {

	app := testutils.NewTestApplication(t)

	arouter := app.InitRouter()

	ts := testutils.NewTestServer(t, arouter)
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snip/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snip/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snip/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snip/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snip/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snip/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.Get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {

	app := testutils.NewTestApplication(t)

	ts := testutils.NewTestServer(t, app.InitRouter())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/signup")
	csrfToken := testutils.ExtractCSRFToken(t, body)

	t.Logf("CSRF token is: %q", csrfToken)
}
