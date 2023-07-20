package handler

import (
	"database/sql"
	"net/http"

	"github.com/Festivals-App/festivals-identity-server/server/status"
)

func GetVersion(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	respondString(w, http.StatusOK, status.VersionString())
}

func GetInfo(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	respondJSON(w, http.StatusOK, status.InfoString())
}

func GetHealth(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	respondCode(w, status.HealthStatus())
}
