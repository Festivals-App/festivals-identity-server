package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/Festivals-App/festivals-identity-server/server/config"
	"github.com/rs/zerolog/log"
)

func GetLog(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	l, err := Log("/var/log/festivals-identity-server/info.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get log")
		respondError(w, http.StatusBadRequest, "Failed to get log")
		return
	}
	respondString(w, http.StatusOK, l)
}

func GetTraceLog(conf *config.Config, w http.ResponseWriter, r *http.Request) {

	l, err := Log("/var/log/festivals-identity-server/trace.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get trace log")
		respondError(w, http.StatusBadRequest, "Failed to get trace log")
		return
	}
	respondString(w, http.StatusOK, l)
}

func Log(location string) (string, error) {

	l, err := os.ReadFile(location)
	if err != nil {
		return "", errors.New("Failed to read log file at: '" + location + "' with error: " + err.Error())
	}
	return string(l), nil
}