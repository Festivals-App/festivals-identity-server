package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
)

func Signup(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	user, err := GetUser(r)
	if err != nil {
		unauthorizedResponse(w)
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func GetUser(r *http.Request) (*User, error) {

	body, readBodyErr := io.ReadAll(r.Body)
	if readBodyErr != nil {
		return nil, readBodyErr
	}
	var userObject User
	err := json.Unmarshal(body, &userObject)
	if err != nil {
		return nil, err
	}
	return &userObject, nil
}
