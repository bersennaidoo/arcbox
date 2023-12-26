package handlers

import (
	"bytes"
	"fmt"
	"net/http"
)

func (h *SnipHandler) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.serverError(w, err)
		return
	}
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
