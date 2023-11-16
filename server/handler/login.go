package handler

import (
	"database/sql"
	"net/http"

	"github.com/Festivals-App/festivals-identity-server/server/database"
	"github.com/Festivals-App/festivals-identity-server/server/model"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Extract the username and password from the request
	// Authorization header. If no Authentication header is present
	// or the header value is invalid, then the 'ok' return value
	// will be false.
	email, password, ok := r.BasicAuth()

	if ok {

		// retrieve user for the given username
		requestedUser, err := GetUserByEmail(db, email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch user.")
			unauthorizedResponse(w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(requestedUser.Password), []byte(password))

		// If the password is correct return the authentication jwt token
		if err == nil {

			servertools.RespondJSON(w, http.StatusOK, requestedUser)
			return
		}
	}

	// If the Authentication header is not present, is invalid, or the
	// username or password is wrong, then set a WWW-Authenticate
	// header to inform the client that we expect them to use basic
	// authentication and send a 401 Unauthorized response.
	unauthorizedResponse(w)
}

func GetUserByEmail(db *sql.DB, email string) (*model.User, error) {

	query := "SELECT * FROM users WHERE `user_email`=?;"
	vars := []interface{}{email}

	rows, err := database.ExecuteRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}
	rows.Next()
	user, err := UserScan(rows)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UserScan(rs *sql.Rows) (model.User, error) {
	var u model.User
	return u, rs.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreateDate, &u.UpdateDate, &u.Role)
}
