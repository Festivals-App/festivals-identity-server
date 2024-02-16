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
	"golang.org/x/crypto/bcrypt"
)

func Signup(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	var signupVars map[string]string
	err = json.Unmarshal(body, &signupVars)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal request body.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	email := signupVars["email"]
	password := signupVars["password"]

	if validEmail(email) && validPassword(password) {

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate password hash from provided password.")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		_, err = database.CreateUserWithEmailAndPasswordHash(db, email, string(passwordHash))
		if err != nil {
			log.Error().Err(err).Msg("Failed to create user with given email and password.")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		servertools.RespondCode(w, http.StatusCreated)
		return
	}

	// If the Authentication header is not present, is invalid, or the username or password is wrong
	servertools.UnauthorizedResponse(w)
}

func Login(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Extract the username and password from the request
	// Authorization header. If no Authentication header is present
	// or the header value is invalid, then the 'ok' return value
	// will be false.
	email, password, ok := r.BasicAuth()

	if ok {

		// retrieve user for the given username
		requestedUser, err := database.GetUserByEmail(db, email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch user.")
			// do i need to mitigate timing attacks on email guessing?
			servertools.UnauthorizedResponse(w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(requestedUser.PasswordHash), []byte(password))
		// If the password is correct return the authentication jwt token
		if err == nil {
			token, err := database.GenerateAccessToken(requestedUser, db, auth)
			if err != nil {
				log.Error().Err(err).Msg("Failed to generate access token for user.")
				servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			servertools.RespondString(w, http.StatusOK, token)
			return
		} else {
			log.Error().Err(err).Msg("The password provided was wrong.")
		}
	}

	// If the Authentication header is not present, is invalid, or the username or password is wrong
	servertools.UnauthorizedResponse(w)
}

func GetUsers(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get users.")
		servertools.UnauthorizedResponse(w)
		return
	}

	users, err := database.GetAllUsers(db)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all users.")
		servertools.UnauthorizedResponse(w)
		return
	}
	servertools.RespondJSON(w, http.StatusOK, users)
}

func ChangePassword(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	if claims.UserID == userID {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read request body.")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		var passwordChangeVars map[string]string
		err = json.Unmarshal(body, &passwordChangeVars)
		if err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal request body.")
			servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		oldpassword := passwordChangeVars["old-password"]
		newpassword := passwordChangeVars["new-password"]

		if oldpassword != newpassword && validPassword(newpassword) && validPassword(oldpassword) {

			// retrieve user for the given username
			requestedUser, err := database.GetUserByID(db, userID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch user.")
				servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(requestedUser.PasswordHash), []byte(oldpassword))
			if err != nil {
				log.Error().Err(err).Msg("Old password is incorrect.")
				servertools.UnauthorizedResponse(w)
				return
			}

			passwordHash, err := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost)
			if err != nil {
				log.Error().Err(err).Msg("Failed to generate password hash from provided password.")
				servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}

			_, err = database.SetPasswordForUser(db, userID, string(passwordHash))
			if err != nil {
				log.Error().Err(err).Msg("Failed to set new password for user.")
				servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}

			servertools.RespondCode(w, http.StatusOK)
			return
		}
	}

	log.Error().Msg("User is not authorized to change password for user.")
	servertools.UnauthorizedResponse(w)
}

func SuspendUser(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get users.")
		servertools.UnauthorizedResponse(w)
		return
	}
	// /users/{objectID}/suspend
	servertools.UnauthorizedResponse(w)
}

func SetUserRole(auth *token.AuthService, claims *token.UserClaims, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	if claims.UserRole != token.ADMIN {
		log.Error().Msg("User is not authorized to get users.")
		servertools.UnauthorizedResponse(w)
		return
	}

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	resourceIDstring, err := resourceID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	resourceID, err := strconv.ParseInt(resourceIDstring, 10, 64)
	if err != nil || (resourceID != int64(token.ADMIN) && resourceID != int64(token.CREATOR)) {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	_, err = database.SetRoleForUser(db, userID, int(resourceID))
	if err != nil {
		log.Error().Err(err).Msg("Failed to set new role for user.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	servertools.RespondCode(w, http.StatusOK)
}

// Set associated resources

func SetFestivalForUser(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	resourceID, err := resourceID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	_, err = database.SetFestivalForUser(db, resourceID, userID)
	if err != nil {
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondCode(w, http.StatusOK)
}

func SetArtistForUser(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	resourceID, err := resourceID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	_, err = database.SetArtistForUser(db, resourceID, userID)
	if err != nil {
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondCode(w, http.StatusOK)
}

func SetLocationForUser(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	resourceID, err := resourceID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	_, err = database.SetLocationForUser(db, resourceID, userID)
	if err != nil {
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondCode(w, http.StatusOK)
}

func RemoveFestivalForUser(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	resourceID, err := resourceID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	_, err = database.RemoveFestivalForUser(db, resourceID, userID)
	if err != nil {
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondCode(w, http.StatusOK)
}

func RemoveArtistForUser(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	resourceID, err := resourceID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	_, err = database.RemoveArtistForUser(db, resourceID, userID)
	if err != nil {
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondCode(w, http.StatusOK)
}

func RemoveLocationForUser(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	userID, err := objectID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	resourceID, err := resourceID(r)
	if err != nil {
		servertools.RespondError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	_, err = database.RemoveLocationForUser(db, resourceID, userID)
	if err != nil {
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	servertools.RespondCode(w, http.StatusOK)
}
