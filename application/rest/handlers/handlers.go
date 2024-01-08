package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/bersennaidoo/arcbox/domain/models"
	"github.com/bersennaidoo/arcbox/infrastructure/repositories/mysql"
	"github.com/gorilla/mux"
	"github.com/kataras/golog"
)

type SnipHandler struct {
	log             *golog.Logger
	snipsRepository *mysql.SnipsRepository
	templateCache   map[string]*template.Template
}

func New(log *golog.Logger, snipsRepository *mysql.SnipsRepository, templateCache map[string]*template.Template) *SnipHandler {
	return &SnipHandler{
		log:             log,
		snipsRepository: snipsRepository,
		templateCache:   templateCache,
	}
}

func (h *SnipHandler) Home(w http.ResponseWriter, r *http.Request) {

	snips, err := h.snipsRepository.Latest()
	if err != nil {
		h.serverError(w, err)
		return
	}
	data := h.newTemplateData(r)
	data.Snips = snips

	h.render(w, http.StatusOK, "home.tmpl", data)
}

func (h *SnipHandler) SnipView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ids := params["id"]
	id, err := strconv.Atoi(ids)
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}

	snip, err := h.snipsRepository.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.notFound(w)
		} else {
			h.serverError(w, err)
		}
		return
	}

	data := h.newTemplateData(r)
	data.Snip = snip

	h.render(w, http.StatusOK, "view.tmpl", data)
}

func (h *SnipHandler) SnipCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form := snipCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := h.snipsRepository.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		h.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)
}

func (h *SnipHandler) SnipCreate(w http.ResponseWriter, r *http.Request) {

	data := h.newTemplateData(r)

	data.Form = snipCreateForm{
		Expires: 365,
	}

	h.render(w, http.StatusOK, "create.tmpl", data)
}
