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

// FESTIVALS

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

func SetFestivalForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_festival_user(`associated_festival`, `associated_user`) VALUES (?, ?);"
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

func RemoveFestivalForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_festival_user WHERE `associated_festival`=? AND `associated_user`=?;"
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

// ARTISTS

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

func SetArtistForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_artist_user(`associated_artist`, `associated_user`) VALUES (?, ?);"
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

func RemoveArtistForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_artist_user WHERE `associated_artist`=? AND `associated_user`=?;"
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

// LOCATIONS

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

func SetLocationForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_location_user(`associated_location`, `associated_user`) VALUES (?, ?);"
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

func RemoveLocationForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_location_user WHERE `associated_location`=? AND `associated_user`=?;"
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

// EVENTS

func GetEventsForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_event` FROM map_event_user WHERE `associated_user`=?;"
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

func SetEventForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_event_user(`associated_event`, `associated_user`) VALUES (?, ?);"
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

func RemoveEventForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_event_user WHERE `associated_event`=? AND `associated_user`=?;"
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

// LINKS

func GetLinksForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_link` FROM map_link_user WHERE `associated_user`=?;"
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

func SetLinkForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_link_user(`associated_link`, `associated_user`) VALUES (?, ?);"
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

func RemoveLinkForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_link_user WHERE `associated_link`=? AND `associated_user`=?;"
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

// IMAGES

func GetImagesForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_image` FROM map_image_user WHERE `associated_user`=?;"
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

func SetImageForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_image_user(`associated_image`, `associated_user`) VALUES (?, ?);"
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

func RemoveImageForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_image_user WHERE `associated_image`=? AND `associated_user`=?;"
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

// PLACES

func GetPlacesForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_place` FROM map_place_user WHERE `associated_user`=?;"
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

func SetPlaceForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_place_user(`associated_place`, `associated_user`) VALUES (?, ?);"
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

func RemovePlaceForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_place_user WHERE `associated_place`=? AND `associated_user`=?;"
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

// TAGS

func GetTagsForUser(db *sql.DB, userID string) ([]int, error) {

	query := "SELECT `associated_tag` FROM map_tag_user WHERE `associated_user`=?;"
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

func SetTagForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "INSERT INTO map_tag_user(`associated_tag`, `associated_user`) VALUES (?, ?);"
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

func RemoveTagForUser(db *sql.DB, objectID string, userID string) (bool, error) {

	query := "DELETE FROM map_tag_user WHERE `associated_tag`=? AND `associated_user`=?;"
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
