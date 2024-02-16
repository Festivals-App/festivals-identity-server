package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	token "github.com/Festivals-App/festivals-identity-server/jwt"
	"github.com/Festivals-App/festivals-identity-server/server/database"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func GetAPIKeys(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	keys, err := database.GetAllAPIKeys(db)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all API keys.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondJSON(w, http.StatusOK, keys)
}

func AddAPIKey(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to create API keys.")
		servertools.UnauthorizedResponse(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	var apikeyValue map[string]string
	err = json.Unmarshal(body, &apikeyValue)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	newAPIKey := apikeyValue["api_key"]
	if newAPIKey == "" {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newAPIKeyComment := apikeyValue["api_key_comment"]
	if newAPIKeyComment == "" {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	apiKey := token.APIKey{
		ID:      0,
		Key:     newAPIKey,
		Comment: newAPIKeyComment,
	}

	err = database.AddAPIKey(db, apiKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add api key.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondCode(w, http.StatusCreated)
}

func UpdateAPIKey(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to update API keys.")
		servertools.UnauthorizedResponse(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	var apikeyValue map[string]string
	err = json.Unmarshal(body, &apikeyValue)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	newAPIKey := apikeyValue["api_key"]
	if newAPIKey == "" {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newAPIKeyComment := apikeyValue["api_key_comment"]
	if newAPIKeyComment == "" {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	keyIDString, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	id, err := strconv.Atoi(keyIDString)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	apiKey := token.APIKey{
		ID:      id,
		Key:     newAPIKey,
		Comment: newAPIKeyComment,
	}

	err = database.UpdateAPIKey(db, apiKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update api key.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondCode(w, http.StatusOK)
}

func DeleteAPIKey(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to delete API keys.")
		servertools.UnauthorizedResponse(w)
		return
	}

	keyID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	err = database.RemoveAPIKey(db, keyID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete api key.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondCode(w, http.StatusOK)
}
