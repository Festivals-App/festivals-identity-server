package handler

import (
	"database/sql"
	"net/http"

	servertools "github.com/Festivals-App/festivals-server-tools"
)

func Refresh(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	user, err := GetUserFromRequest(r)
	if err != nil {
		unauthorizedResponse(w)
		return
	}
	servertools.RespondJSON(w, http.StatusOK, user)
}
