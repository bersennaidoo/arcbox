package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bersennaidoo/arcbox/domain/models"
	"github.com/bersennaidoo/arcbox/infrastructure/repositories/mysql"
	"github.com/kataras/golog"
)

type SnipHandler struct {
	log             *golog.Logger
	snipsRepository *mysql.SnipsRepository
}

func New(log *golog.Logger, snipsRepository *mysql.SnipsRepository) *SnipHandler {
	return &SnipHandler{
		log:             log,
		snipsRepository: snipsRepository,
	}
}

func (h *SnipHandler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w)
		return
	}

	snips, err := h.snipsRepository.Latest()
	if err != nil {
		h.serverError(w, err)
		return
	}
	for _, snip := range snips {
		fmt.Fprintf(w, "%+v\n", snip)
	}

	/*files := []string{
		"./hci/html/base.tmpl",
		"./hci/html/partials/nav.tmpl",
		"./hci/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.serverError(w, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		h.serverError(w, err)
		http.Error(w, "Internal Server Error", 500)
	}*/
}

func (h *SnipHandler) SnipView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%+v", snip)
}

func (h *SnipHandler) SnipCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := h.snipsRepository.Insert(title, content, expires)
	if err != nil {
		h.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snip/view?id=%d", id), http.StatusSeeOther)
}
