package handlers

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/bersennaidoo/arcbox/domain/models"
	"github.com/bersennaidoo/arcbox/hci"
	"github.com/bersennaidoo/arcbox/transport/http/validator"
	"github.com/justinas/nosurf"
)

type snipCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type templateData struct {
	CurrentYear     int
	Snip            *models.Snip
	Snips           []*models.Snip
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(hci.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(hci.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}

func (h *Handler) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           h.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: h.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}
