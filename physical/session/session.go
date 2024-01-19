package session

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

func New(db *sql.DB) *scs.SessionManager {

	sessionManager := scs.New()
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	return sessionManager
}
