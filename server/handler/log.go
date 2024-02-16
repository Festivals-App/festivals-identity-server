package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	token "github.com/Festivals-App/festivals-identity-server/jwt"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func GetLog(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get log file.")
		servertools.UnauthorizedResponse(w)
		return
	}

	l, err := Log("/var/log/festivals-identity-server/info.log")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get log")
		servertools.RespondError(w, http.StatusBadRequest, "Failed to get log")
		return
	}
	servertools.RespondString(w, http.StatusOK, l)
}

func GetTraceLog(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get trace log file.")
		servertools.UnauthorizedResponse(w)
		return
	}

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
