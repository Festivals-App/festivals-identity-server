package database

import (
	"database/sql"
	"errors"

	token "github.com/Festivals-App/festivals-identity-server/auth"
)

func GetAllServiceKeys(db *sql.DB) ([]token.ServiceKey, error) {

	query := "SELECT * FROM service_keys;"
	vars := []interface{}{}

	rows, err := executeRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}
	keys := []token.ServiceKey{}
	for rows.Next() {
		key, err := serviceKeyScan(rows)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func AddServiceKey(db *sql.DB, key token.ServiceKey) error {

	query := "INSERT INTO service_keys(`service_key`, `service_key_comment`) VALUES (?, ?);"
	vars := []interface{}{key.Key, key.Comment}

	result, err := executeQuery(db, query, vars)
	if err != nil {
		return err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	if insertID == 0 {
		return errors.New("failed to insert new service key without mysql error")
	}

	return nil
}

func UpdateServiceKey(db *sql.DB, key token.ServiceKey) error {

	query := "UPDATE service_keys SET `service_key`=?, `service_key_comment`=? WHERE `service_key_id`=?;"
	vars := []interface{}{key.Key, key.Comment, key.ID}

	result, err := executeQuery(db, query, vars)
	if err != nil {
		return err
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numOfAffectedRows != 1 {
		return errors.New("failed to update service key without mysql error")
	}
	return nil
}

func RemoveServiceKey(db *sql.DB, keyID string) error {

	query := "DELETE FROM service_keys WHERE `service_key_id`=?;"
	vars := []interface{}{keyID}

	result, err := executeQuery(db, query, vars)
	if err != nil {
		return err
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numOfAffectedRows != 1 {
		return errors.New("failed to delete service key without mysql error")
	}
	return nil
}
