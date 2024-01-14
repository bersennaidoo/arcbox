package handlers

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/bersennaidoo/arcbox/application/rest/validator"
	"github.com/bersennaidoo/arcbox/domain/models"
	"github.com/bersennaidoo/arcbox/hci"
)

type snipCreateForm struct {
	Title               string `schema:"title"`
	Content             string `schema:"content"`
	Expires             int    `schema:"expires"`
	validator.Validator `schema:"-"`
}

type userLoginForm struct {
	Email               string `schema:"email"`
	Password            string `schema:"password"`
	validator.Validator `schema:"-"`
}

type templateData struct {
	CurrentYear     int
	Snip            *models.Snip
	Snips           []*models.Snip
	Form            any
	Flash           string
	IsAuthenticated bool
}

type userSignupForm struct {
	Name                string `schema:"name"`
	Email               string `schema:"email"`
	Password            string `schema:"password"`
	validator.Validator `schema:"-"`
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
	}
}
