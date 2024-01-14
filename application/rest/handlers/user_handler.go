package handlers

import (
	"errors"
	"net/http"

	"github.com/bersennaidoo/arcbox/application/rest/validator"
	"github.com/bersennaidoo/arcbox/domain/models"
)

func (h *Handler) UserSignup(w http.ResponseWriter, r *http.Request) {

	data := h.newTemplateData(r)
	data.Form = userSignupForm{}

	h.render(w, http.StatusOK, "signup.tmpl", data)
}

func (h *Handler) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := h.decodePostForm(r, &form)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = h.usersRepository.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := h.newTemplateData(r)
			data.Form = form
			h.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			h.serverError(w, err)
		}
		return
	}

	h.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {

	data := h.newTemplateData(r)
	data.Form = userLoginForm{}

	h.render(w, http.StatusOK, "login.tmpl", data)
}

func (h *Handler) UserLoginPost(w http.ResponseWriter, r *http.Request) {

	var form userLoginForm
	err := h.decodePostForm(r, &form)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := h.usersRepository.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := h.newTemplateData(r)
			data.Form = form
			h.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			h.serverError(w, err)
		}
		return
	}
	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snip/create", http.StatusSeeOther)
}

func (h *Handler) UserLogoutPost(w http.ResponseWriter, r *http.Request) {

	err := h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.sessionManager.Remove(r.Context(), "authenticatedUserID")

	h.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
