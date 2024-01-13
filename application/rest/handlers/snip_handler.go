package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/bersennaidoo/arcbox/application/rest/validator"
	"github.com/bersennaidoo/arcbox/domain/models"
	"github.com/bersennaidoo/arcbox/infrastructure/repositories/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/kataras/golog"
)

type SnipHandler struct {
	log             *golog.Logger
	snipsRepository *mysql.SnipsRepository
	usersRepository *mysql.UsersRepository
	templateCache   map[string]*template.Template
	formDecoder     *schema.Decoder
	sessionManager  *scs.SessionManager
}

func New(log *golog.Logger, snipsRepository *mysql.SnipsRepository, usersRepository *mysql.UsersRepository,
	templateCache map[string]*template.Template, formDecoder *schema.Decoder,
	sessionManager *scs.SessionManager) *SnipHandler {
	return &SnipHandler{
		log:             log,
		snipsRepository: snipsRepository,
		usersRepository: usersRepository,
		templateCache:   templateCache,
		formDecoder:     formDecoder,
		sessionManager:  sessionManager,
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

	var form snipCreateForm

	err := h.decodePostForm(r, &form)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
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

	h.sessionManager.Put(r.Context(), "flash", "Snip successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snip/view/%d", id), http.StatusSeeOther)
}

func (h *SnipHandler) SnipCreate(w http.ResponseWriter, r *http.Request) {

	data := h.newTemplateData(r)

	data.Form = snipCreateForm{
		Expires: 365,
	}

	h.render(w, http.StatusOK, "create.tmpl", data)
}
