package handler

import (
	"net/http"

	"github.com/Festivals-App/festivals-identity-server/server/config"
	"github.com/Festivals-App/festivals-identity-server/server/status"
)

func GetVersion(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	respondString(w, http.StatusOK, status.VersionString())
}

func GetInfo(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	respondJSON(w, http.StatusOK, status.InfoString())
}

func GetHealth(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	respondCode(w, status.HealthStatus())
}
