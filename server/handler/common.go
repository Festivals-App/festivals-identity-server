package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func validEmail(email string) bool {

	if len(email) < 6 {
		return false
	}

	return true
}

func validPassword(password string) bool {

	if len(password) < 8 {
		return false
	}

	return true
}

func objectID(r *http.Request) (string, error) {
	return chi.URLParam(r, "objectID"), nil
}

func resourceID(r *http.Request) (string, error) {
	return chi.URLParam(r, "resourceID"), nil
}
