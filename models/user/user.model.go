package userModel

import (
	"database/sql"
	"olx-clone/infra/db"
)

func InsertUser(body UserBody, hashedPassword string) error {
	database := db.Mgr.DBConn
	query := `INSERT INTO users(name, email, password, number) VALUES($1, $2, $3, $4)`
	_, err := database.Exec(query, body.Name, body.Email, hashedPassword, body.Number)
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
	query := `SELECT id, name, email, number FROM users ORDER BY id LIMIT $1 OFFSET $2`
	return database.Query(query, itemsPerPage, offset)
}
