package userModel

import (
	"context"
	"database/sql"
	"errors"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	"olx-clone/infra/db"

	"golang.org/x/crypto/bcrypt"
)

var database = db.Mgr.DBConn

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func UserAlreadyExistsWithUsername(username string) bool {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"

	err := database.QueryRow(query, username).Scan(&exists)

	if err != nil {
		logger.Log.Println(err)
		return false
	}

	return exists
}

func GetUserByUsername(username string) (User, error) {
	var userModel User
	validUserName := general.ValidUserName(username)
	if !validUserName {
		return userModel, errors.New("invalid username")
	}

	query := "SELECT * FROM users WHERE username = $1"
	err := database.GetContext(context.TODO(), &userModel, query, username)
	if err == nil {
		return userModel, nil
	}
	return userModel, err
}

func UpdateUser(context context.Context, username string, body UserUpdateBody) error {
	query := "UPDATE users SET username = $2, email = $3, phone = $4, avatar = $5, location = $6, coordinates = $7, fullname = $8 WHERE username = $1"
	res, err := database.ExecContext(context, query, username, body.Username, body.Email, body.Phone, body.Avatar, body.Location, body.Coordinates, body.Fullname)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("could not update user")
	}

	return nil
}

func IsValidUser(username, password string) (User, error) {
	var userModel User
	validUserName := general.ValidUserName(username)

	if !validUserName {
		return userModel, errors.New("invalid username or password")
	}

	user, err := GetUserByUsername(username)
	if err != nil {
		return userModel, err
	}

	valid := checkPasswordHash(password, user.Password)

	if !valid {
		return userModel, errors.New("invalid username or password")
	}

	return userModel, nil
}

func DeleteUserByUsername(username string) error {
	query := "DELETE FROM users WHERE username = $1"
	res, err := database.ExecContext(context.TODO(), query, username)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user already deleted or user not found")
	}

	return err
}

func InsertUser(body UserBody, password string) error {
	query := `INSERT INTO users(username, fullname, email, password, phone, avatar, location, coordinates) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := database.Exec(query, body.Username, body.Fullname, body.Email, password, body.Phone, body.Avatar, body.Location, body.Coordinates)
	if err != nil {
		return err
	}
	return nil
}

func GetUsersListPaginatedValue(itemsPerPage, offset int) (*sql.Rows, error) {
	query := `SELECT id, username, email, phone FROM users ORDER BY id LIMIT $1 OFFSET $2`
	return database.Query(query, itemsPerPage, offset)
}
