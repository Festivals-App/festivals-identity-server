package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Select(db *sql.DB, table string, objectIDs []int) (*sql.Rows, error) {

	var query string
	var vars []int
	if len(objectIDs) == 0 {
		query = "SELECT * FROM " + table + "s;"
		vars = []int{}
	} else {
		placeholder := DBPlaceholderForIDs(objectIDs)
		query = "SELECT * FROM " + table + "s WHERE " + table + "_id IN (" + placeholder + ");"
		vars = objectIDs
	}
	return ExecuteRowQuery(db, query, InterfaceInt(vars))
}

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

func ExecuteQuery(db *sql.DB, query string, args []interface{}) (sql.Result, error) {

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
