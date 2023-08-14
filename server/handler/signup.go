package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Festivals-App/festivals-identity-server/server/model"
)

func Signup(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	user, err := GetUserFromRequest(r)
	if err != nil {
		unauthorizedResponse(w)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func GetUserFromRequest(r *http.Request) (*model.User, error) {

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		return nil, readBodyErr
	}
	var userObject model.User
	err := json.Unmarshal(body, &userObject)
	if err != nil {
		return nil, err
	}
	return &userObject, nil
}
