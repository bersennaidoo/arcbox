package handlers_test

import (
	"net/http"
	"testing"

	"github.com/bersennaidoo/arcbox/foundation/assert"
	"github.com/bersennaidoo/arcbox/foundation/testutils"
)

func TestPing(t *testing.T) {

	app := testutils.NewTestApplication(t)

	pingr, _, _, _ := app.InitRouter()

	ts := testutils.NewTestServer(t, pingr)
	defer ts.Close()

	code, _, body := ts.Get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)

	assert.Equal(t, body, "OK")
}
