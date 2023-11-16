package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func GetLog(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	l, err := Log("/var/log/festivals-identity-server/info.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get log")
		servertools.RespondError(w, http.StatusBadRequest, "Failed to get log")
		return
	}
	servertools.RespondString(w, http.StatusOK, l)
}

func GetTraceLog(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	l, err := Log("/var/log/festivals-identity-server/trace.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get trace log")
		servertools.RespondError(w, http.StatusBadRequest, "Failed to get trace log")
		return
	}
	servertools.RespondString(w, http.StatusOK, l)
}

func Log(location string) (string, error) {

	l, err := os.ReadFile(location)
	if err != nil {
		return "", errors.New("Failed to read log file at: '" + location + "' with error: " + err.Error())
	}
	return string(l), nil
}
