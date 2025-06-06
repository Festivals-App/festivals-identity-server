package database

import (
	"database/sql"

	token "github.com/Festivals-App/festivals-identity-server/auth"
)

type Entity string

const (
	Festival Entity = "festival"
	Artist   Entity = "artist"
	Location Entity = "location"
	Event    Entity = "event"
	Link     Entity = "link"
	Image    Entity = "image"
	Place    Entity = "place"
	Tag      Entity = "tag"
)

func executeRowQuery(db *sql.DB, query string, args []interface{}) (*sql.Rows, error) {

	rows, err := db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rows, nil
}

func executeQuery(db *sql.DB, query string, args []interface{}) (sql.Result, error) {

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	err = stmt.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func userScan(rs *sql.Rows) (token.User, error) {
	var u token.User
	return u, rs.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreateDate, &u.UpdateDate, &u.Role)
}

func userSummaryScan(rs *sql.Rows) (token.UserSummary, error) {
	var u token.UserSummary
	return u, rs.Scan(&u.ID, &u.Email, &u.CreateDate, &u.UpdateDate, &u.Role)
}

func apiKeyScan(rs *sql.Rows) (token.APIKey, error) {
	var u token.APIKey
	return u, rs.Scan(&u.ID, &u.Key, &u.Comment)
}

func serviceKeyScan(rs *sql.Rows) (token.ServiceKey, error) {
	var u token.ServiceKey
	return u, rs.Scan(&u.ID, &u.Key, &u.Comment)
}
