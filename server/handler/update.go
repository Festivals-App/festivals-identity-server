package handler

import (
	"database/sql"
	"net/http"

	"github.com/Festivals-App/festivals-identity-server/server/status"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func MakeUpdate(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	newVersion, err := servertools.RunUpdate(status.ServerVersion, "Festivals-App", "festivals-identity-server", "/usr/local/festivals-identity-server/update.sh")
	if err != nil {
		log.Error().Err(err).Msg("Failed to update")
		servertools.RespondError(w, http.StatusInternalServerError, "Failed to update")
		return
	}
	servertools.RespondString(w, http.StatusAccepted, newVersion)
}
