package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/kataras/golog"
)

type SnipHandler struct {
	log *golog.Logger
}

func New(log *golog.Logger) *SnipHandler {
	return &SnipHandler{
		log: log,
	}
}

func (h *SnipHandler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.notFound(w)
		return
	}

	files := []string{
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
	}
}

func (h *SnipHandler) SnipView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		h.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specific snip with ID %d...", id)
}

func (h *SnipHandler) SnipCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		h.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
