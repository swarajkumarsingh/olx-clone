package userModel

import (
	"database/sql"
	"olx-clone/infra/db"
)

var database = db.Mgr.DBConn

func UserAlreadyExists(username string) bool {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = '$1')"
	err := database.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func InsertUser(body UserBody, password string) error {

	query := `INSERT INTO users(username, fullname, email, password, phone, avatar, location, coordinates) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := database.Exec(query, body.Username, body.Fullname, body.Email, password, body.Number, body.Avatar, body.Location, body.Coordinates)
	if err != nil {
		return err
	}
	return nil
}

func GetUsersListPaginatedValue(itemsPerPage int, offset int) (*sql.Rows, error) {
	query := `SELECT id, username, email, phone FROM users ORDER BY id LIMIT $1 OFFSET $2`
	return database.Query(query, itemsPerPage, offset)
}
