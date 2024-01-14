package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

func (h *Handler) render(w http.ResponseWriter, status int, page string, data *templateData) {
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

func (h *Handler) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = h.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *schema.ConversionError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (h *Handler) isAuthenticated(r *http.Request) bool {
	return h.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
