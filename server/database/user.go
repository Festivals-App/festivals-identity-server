package database

import (
	"database/sql"

	token "github.com/Festivals-App/festivals-identity-server/jwt"
)

func GetAllUsers(db *sql.DB) ([]*token.User, error) {

	query := "SELECT * FROM users;"
	vars := []interface{}{}

	rows, err := executeRowQuery(db, query, vars)
	if err != nil {
		return nil, err
	}

	keys := []*token.User{}
	for rows.Next() {
		key, err := userScan(rows)
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

func GetFestivalsForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_festival` FROM map_festival_user WHERE `associated_user`=?;"
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

func SetFestivalForUser(db *sql.DB, festivalID string, userID string) (bool, error) {

	query := "INSERT INTO map_festival_user(`associated_festival`, `associated_user`) VALUES (?, ?);"
	vars := []interface{}{festivalID, userID}
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

func RemoveFestivalForUser(db *sql.DB, festivalID string, userID string) (bool, error) {

	query := "DELETE FROM map_festival_user WHERE `associated_festival`=? AND `associated_user`=?;"
	vars := []interface{}{festivalID, userID}

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

func GetArtistsForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_artist` FROM map_artist_user WHERE `associated_user`=?;"
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

func SetArtistForUser(db *sql.DB, artistID string, userID string) (bool, error) {

	query := "INSERT INTO map_artist_user(`associated_artist`, `associated_user`) VALUES (?, ?);"
	vars := []interface{}{artistID, userID}

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

func RemoveArtistForUser(db *sql.DB, arrtistID string, userID string) (bool, error) {

	query := "DELETE FROM map_artist_user WHERE `associated_artist`=? AND `associated_user`=?;"
	vars := []interface{}{arrtistID, userID}

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

func GetLocationsForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_location` FROM map_location_user WHERE `associated_user`=?;"
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

func SetLocationForUser(db *sql.DB, locationID string, userID string) (bool, error) {

	query := "INSERT INTO map_location_user(`associated_location`, `associated_user`) VALUES (?, ?);"
	vars := []interface{}{locationID, userID}

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

func RemoveLocationForUser(db *sql.DB, locationID string, userID string) (bool, error) {

	query := "DELETE FROM map_location_user WHERE `associated_location`=? AND `associated_user`=?;"
	vars := []interface{}{locationID, userID}

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
