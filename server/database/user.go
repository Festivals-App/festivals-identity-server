package database

import (
	"database/sql"

	token "github.com/Festivals-App/festivals-identity-server/auth"
)

func GetAllUserSummaries(db *sql.DB) ([]*token.UserSummary, error) {

	query := "SELECT user_id, user_email, user_createdat, user_updatedat, user_role FROM users;"
	vars := []any{}

	rows, err := executeRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}

	keys := []*token.UserSummary{}
	for rows.Next() {
		key, err := userSummaryScan(rows)
		if err != nil {
			return nil, err
		}
		keys = append(keys, &key)
	}
	return keys, nil
}

func GetUserByEmail(db *sql.DB, email string) (*token.User, error) {

	query := "SELECT * FROM users WHERE `user_email`=?;"
	vars := []interface{}{email}

	rows, err := executeRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}
	rows.Next()
	user, err := userScan(rows)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(db *sql.DB, userID string) (*token.User, error) {

	query := "SELECT * FROM users WHERE `user_id`=?;"
	vars := []interface{}{userID}

	rows, err := executeRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}
	rows.Next()
	user, err := userScan(rows)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUserWithEmailAndPasswordHash(db *sql.DB, email string, passwordhash string) (bool, error) {

	query := "INSERT INTO `users`(`user_email`, `user_password`, `user_role`) VALUES (?, ?, ?);"
	vars := []interface{}{email, passwordhash, token.CREATOR}

	result, err := executeQuery(db, query, vars)
	if err != nil {
		return false, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return false, err
	}

	return insertID != 0, nil
}

func SetPasswordForUser(db *sql.DB, userID string, newpasswordhash string) (bool, error) {

	query := "UPDATE `users` SET `user_password`=? WHERE `user_id`=?;"
	vars := []interface{}{newpasswordhash, userID}

	result, err := executeQuery(db, query, vars)
	if err != nil {
		return false, err
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if numOfAffectedRows != 1 {
		return false, err
	}
	return true, nil
}

func SetRoleForUser(db *sql.DB, userID string, newUserRole int) (bool, error) {

	query := "UPDATE `users` SET `user_role`=? WHERE `user_id`=?;"
	vars := []interface{}{newUserRole, userID}

	result, err := executeQuery(db, query, vars)
	if err != nil {
		return false, err
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if numOfAffectedRows != 1 {
		return false, err
	}
	return true, nil
}

func GetEntitiesForUser(entity Entity, db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_" + string(entity) + "` FROM map_" + string(entity) + "_user WHERE `associated_user`=?;"
	vars := []interface{}{userID}

	rows, err := executeRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}
	ids := []int{}
	for rows.Next() {
		var fid int
		err = rows.Scan(&fid)
		if err != nil {
			return nil, err
		}

		ids = append(ids, fid)
	}
	return ids, nil
}

func SetEntityForUser(entity Entity, db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_" + string(entity) + "_user(`associated_" + string(entity) + "`, `associated_user`) VALUES (?, ?);"
	vars := []interface{}{objectID, userID}
	result, err := executeQuery(db, query, vars)
	if err != nil {
		return false, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return false, err
	}
	return insertID != 0, nil
}

func RemoveEntityForUser(entity Entity, db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_" + string(entity) + "_user WHERE `associated_" + string(entity) + "`=? AND `associated_user`=?;"
	vars := []interface{}{objectID, userID}

	result, err := executeQuery(db, query, vars)
	if err != nil {
		return false, err
	}
	numOfAffectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if numOfAffectedRows != 1 {
		return false, err
	}
	return true, nil
}
