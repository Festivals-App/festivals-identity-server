package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ExecuteRowQuery(db *sql.DB, query string, args []interface{}) (*sql.Rows, error) {

	rows, err := db.Query(query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rows, nil
}
