package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	token "github.com/Festivals-App/festivals-identity-server/auth"
	"github.com/Festivals-App/festivals-identity-server/server/database"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func GetServiceKeys(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	keys, err := database.GetAllServiceKeys(db)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all service keys.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondJSON(w, http.StatusOK, keys)
}

func AddServiceKey(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to create service keys.")
		servertools.UnauthorizedResponse(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	var servicekeyValue map[string]string
	err = json.Unmarshal(body, &servicekeyValue)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	newServiceKey := servicekeyValue["service_key"]
	if newServiceKey == "" {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	newServiceKeyComment := servicekeyValue["service_key_comment"]
	if newServiceKeyComment == "" {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	serviceKey := token.ServiceKey{
		ID:      0,
		Key:     newServiceKey,
		Comment: newServiceKeyComment,
	}

	err = database.AddServiceKey(db, serviceKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add service key.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondCode(w, http.StatusCreated)
}

func UpdateServiceKey(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to update service keys.")
		servertools.UnauthorizedResponse(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	var servicekeyValue map[string]string
	err = json.Unmarshal(body, &servicekeyValue)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	serviceKey := servicekeyValue["service_key"]
	if serviceKey == "" {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	serviceKeyComment := servicekeyValue["service_key_comment"]
	if serviceKeyComment == "" {
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

	key := token.ServiceKey{
		ID:      id,
		Key:     serviceKey,
		Comment: serviceKeyComment,
	}

	err = database.UpdateServiceKey(db, key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update service key.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondCode(w, http.StatusOK)
}

func DeleteServiceKey(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to delete service keys.")
		servertools.UnauthorizedResponse(w)
		return
	}

	keyID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	err = database.RemoveServiceKey(db, keyID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete api key.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondCode(w, http.StatusOK)
}
