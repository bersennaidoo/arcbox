package handlers

import (
	"fmt"
	"net/http"
)

func (h *SnipHandler) UserSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Displaying a HTML form for signing up a new user...")
}

func (h *SnipHandler) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (h *SnipHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (h *SnipHandler) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (h *SnipHandler) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
