package userModel

import (
	"context"
	"database/sql"
	"errors"
	"olx-clone/constants"
	"olx-clone/constants/messages"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	"olx-clone/infra/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var database = db.Mgr.DBConn

func GetOtpFromDB(context context.Context, username string) (string, error) {
	var OTPSecret string
	var OTPExpiration time.Time

	query := "SELECT otp, otp_expiration FROM users WHERE username = $1"
	err := database.QueryRowContext(context, query, username).Scan(&OTPSecret, &OTPExpiration)
	if err != nil {
		return OTPSecret, err
	}

	if time.Now().After(OTPExpiration) {
		return "", errors.New(messages.OTPExpiredMessage)
	}

	return OTPSecret, nil
}

func ResetOtpAndOtpExpiration(context context.Context, username string) error {
	_, err := database.ExecContext(context, "UPDATE users SET otp = '', otp_expiration = NULL WHERE username = $1", username)
	if err != nil {
		return err
	}
	return nil
}

func GetViewedProducts(userId string) (any, error) {
	var model ViewedProductStruct

	query := "SELECT * FROM product_views WHERE user_id = $1"
	err := database.GetContext(context.TODO(), &model, query, userId)
	if err == nil {
		return model, nil
	}
	return model, err
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

func GetUserByUsername(context context.Context, username string) (User, error) {
	var userModel User
	validUserName := general.ValidUserName(username)
	if !validUserName {
		return userModel, errors.New("invalid username")
	}

	query := "SELECT * FROM users WHERE username = $1"
	err := database.GetContext(context, &userModel, query, username)
	if err == nil {
		return userModel, nil
	}
	return userModel, err
}

func GetUserByUsernameWithUserId(context context.Context, userId string) (User, error) {
	var userModel User
	validUserName := general.ValidUserName(userId)
	if !validUserName {
		return userModel, errors.New("invalid username")
	}

	query := "SELECT * FROM users WHERE id = $1"
	err := database.GetContext(context, &userModel, query, userId)
	if err == nil {
		return userModel, nil
	}
	return userModel, err
}

func SaveOTPAndExpirationInDB(context context.Context, username, otp string, expiration any) error {
	query := "UPDATE users SET otp = $2, otp_expiration = $3 WHERE username = $1"
	res, err := database.ExecContext(context, query, username, otp, expiration)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("could not update user")
	}

	return nil
}

func UpdateUser(context context.Context, username string, body UserUpdateBody) error {
	query := "UPDATE users SET username = $2, avatar = $3, location = $4, coordinates = $5, fullname = $6 WHERE username = $1"
	res, err := database.ExecContext(context, query, username, body.Username, body.Avatar, body.Location, body.Coordinates, body.Fullname)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("could not update user")
	}

	return nil
}

func CheckIfUsernameExists(context context.Context, username string) (User, error) {
	var user User
	user, err := GetUserByUsername(context, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New(messages.UserNotFoundMessage)
		}
		return user, err
	}

	return user, nil
}

func CheckIfUsernameExistsWithId(context context.Context, userId string) (User, error) {
	var user User
	user, err := GetUserByUsernameWithUserId(context, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New(messages.UserNotFoundMessage)
		}
		return user, err
	}

	return user, nil
}

func CheckIfCurrentPasswordIsValid(context context.Context, username, password string) error {
	valid := general.ValidUserName(username)
	if !valid {
		return errors.New(messages.InvalidCredentialsMessage)
	}

	user, err := GetUserByUsername(context, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user does not exists")
		}
		return err
	}

	valid = checkPasswordHash(password, user.Password)
	if !valid {
		return errors.New(messages.InvalidCredentialsMessage)
	}

	return nil
}

func UpdatePassword(context context.Context, username, newPassword string) error {
	password, err := hashPassword(newPassword)
	if err != nil {
		return errors.New(messages.SomethingWentWrongMessage)
	}

	query := "UPDATE users SET password = $2 WHERE username = $1"
	res, err := database.ExecContext(context, query, username, password)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("could not update user")
	}

	return nil
}

func IsValidUser(context context.Context, username, password string) (User, error) {
	var userModel User
	validUserName := general.ValidUserName(username)

	if !validUserName {
		return userModel, errors.New(messages.InvalidCredentialsMessage)
	}

	user, err := GetUserByUsername(context, username)
	if err != nil {
		return userModel, err
	}

	valid := checkPasswordHash(password, user.Password)

	if !valid {
		return userModel, errors.New(messages.InvalidCredentialsMessage)
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

func GetViewedProductsListPaginatedValue(userId string, itemsPerPage, offset int) (*sql.Rows, error) {
	// query := `SELECT id, username, email, phone FROM users ORDER BY id LIMIT $1 OFFSET $2`
	query2 := `SELECT id, user_id, product_id FROM product_views WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	return database.Query(query2, userId, itemsPerPage, offset)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptHashingCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
