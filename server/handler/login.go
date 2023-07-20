package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

/*
		`user_id` 			    int unsigned 	 	NOT NULL AUTO_INCREMENT 											    COMMENT 'The id of the user.',
		`user_name` 	  	    varchar(225) 		NOT NULL 												                COMMENT 'The name of the user. The name needs to be unique.',
		`user_email` 		    varchar(255)		NOT NULL													            COMMENT 'The email of the user. The email needs to be unique.',
		`user_password` 	    varchar(225) 	  	NOT NULL 												                COMMENT '',
		`user_tokenhash` 		varchar(15) 	  	NOT NULL 											                    COMMENT '',
		`user_createdat` 		timestamp 			NOT NULL DEFAULT current_timestamp()					      		    COMMENT '',
		`user_updatedat` 		timestamp 			NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()	    COMMENT '',
	    `user_role` 	  	    tinyint 		    NOT NULL DEFAULT 0											            COMMENT 'The role of the user.',
*/

type User struct {
	ID         int         `json:"user_id" sql:"user_id"`
	Name       string      `json:"user_name" sql:"user_name"`
	Email      string      `json:"user_email" sql:"user_email"`
	Password   string      `json:"user_password" sql:"user_password"`
	CreateDate time.Time   `json:"user_createdat" sql:"user_createdat"`
	UpdateDate time.Time   `json:"user_updatedat" sql:"user_updatedat"`
	Role       int         `json:"user_role" sql:"user_role"`
	Include    interface{} `json:"include,omitempty"`
}

func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// Extract the username and password from the request
	// Authorization header. If no Authentication header is present
	// or the header value is invalid, then the 'ok' return value
	// will be false.
	email, password, ok := r.BasicAuth()

	if ok {

		// retrieve user for the given username
		requestedUser, err := UserWithEmail(db, email)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch user.")
			unauthorizedResponse(w)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(requestedUser.Password), []byte(password))

		// If the password is correct return the authentication jwt token
		if err == nil {
			respondJSON(w, http.StatusOK, requestedUser)
			return
		}
	}

	// If the Authentication header is not present, is invalid, or the
	// username or password is wrong, then set a WWW-Authenticate
	// header to inform the client that we expect them to use basic
	// authentication and send a 401 Unauthorized response.
	unauthorizedResponse(w)
}

func UserWithEmail(db *sql.DB, email string) (*User, error) {

	query := "SELECT * FROM users WHERE `user_email`=?;"
	vars := []interface{}{email}

	rows, err := db.Query(query, vars...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user, err := UserScan(rows)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserScan(rs *sql.Rows) (User, error) {
	var u User
	return u, rs.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreateDate, &u.UpdateDate, &u.Role, &u.Include)
}
