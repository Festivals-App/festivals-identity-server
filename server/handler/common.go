package handler

import (
	"net/http"
)

func unauthorizedResponse(w http.ResponseWriter) {

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
