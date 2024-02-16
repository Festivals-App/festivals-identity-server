package database

import (
	"database/sql"

	token "github.com/Festivals-App/festivals-identity-server/jwt"
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

func apiKeyScan(rs *sql.Rows) (token.APIKey, error) {
	var u token.APIKey
	return u, rs.Scan(&u.ID, &u.Key, &u.Comment)
}

func serviceKeyScan(rs *sql.Rows) (token.ServiceKey, error) {
	var u token.ServiceKey
	return u, rs.Scan(&u.ID, &u.Key, &u.Comment)
}
