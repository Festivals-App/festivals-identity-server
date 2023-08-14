package handler

import (
	"database/sql"
	"net/http"
)

func Refresh(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	user, err := GetUserFromRequest(r)
	if err != nil {
		unauthorizedResponse(w)
		return
	}
	respondJSON(w, http.StatusOK, user)
}
