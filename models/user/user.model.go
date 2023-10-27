package userModel

import (
	"database/sql"
	"olx-clone/infra/db"
)

func InsertUser(body UserBody, password string) error {
	database := db.Mgr.DBConn
	query := `INSERT INTO users(username, fullname, email, password, phone) VALUES($1, $2, $3, $4, $5)`
	_, err := database.Exec(query, body.Username, body.Fullname, body.Email, password, body.Number)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsersQuery(itemsPerPage int, offset int) (*sql.Rows, error) {
	database := db.Mgr.DBConn
	query := `SELECT id, name, email FROM users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := database.Query(query, itemsPerPage, offset)
	if err != nil {
		return rows, err
	}
	defer rows.Close()
	return rows, nil
}

func GetUsersListPaginatedValue(itemsPerPage int, offset int) (*sql.Rows, error) {
	database := db.Mgr.DBConn
	query := `SELECT id, username, email, phone FROM users ORDER BY id LIMIT $1 OFFSET $2`
	return database.Query(query, itemsPerPage, offset)
}
