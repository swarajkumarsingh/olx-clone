package sellerModel

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

func CreateSeller(body SellerBody, password string) error {
	if body.Avatar == "" {
		query := `INSERT INTO sellers(username, fullname, email, password, phone, location, coordinates, description) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := database.Exec(query, body.Username, body.Fullname, body.Email, password, body.Phone, body.Location, body.Coordinates, body.Description)
		if err != nil {
			return err
		}
		return nil
	}

	query := `INSERT INTO sellers(username, fullname, email, password, phone, avatar, location, coordinates, description) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := database.Exec(query, body.Username, body.Fullname, body.Email, password, body.Phone, body.Avatar, body.Location, body.Coordinates, body.Description)
	if err != nil {
		return err
	}
	return nil
}

func SellerAlreadyExistsWithUsername(username string) bool {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM sellers WHERE username = $1)"

	err := database.QueryRow(query, username).Scan(&exists)

	if err != nil {
		logger.Log.Println(err)
		return false
	}

	return exists
}

func GetSellerListPaginatedValue(itemsPerPage, offset int) (*sql.Rows, error) {
	query := `SELECT id, username, email, phone FROM sellers ORDER BY id LIMIT $1 OFFSET $2`
	return database.Query(query, itemsPerPage, offset)
}

func GetSellerByUsername(context context.Context, username string) (Seller, error) {
	var userModel Seller
	validUserName := general.ValidUserName(username)
	if !validUserName {
		return userModel, errors.New("invalid username")
	}

	query := "SELECT * FROM sellers WHERE username = $1"
	err := database.GetContext(context, &userModel, query, username)
	if err == nil {
		return userModel, nil
	}
	return userModel, err
}

func UpdateSeller(context context.Context, username string, body SellerUpdateStruct) error {
	query := "UPDATE sellers SET username = $2, avatar = $3, location = $4, coordinates = $5, fullname = $6, description = $7, city = $8, state = $9, country = $10, zip_code = $11 WHERE username = $1"
	res, err := database.ExecContext(context, query, username, body.Username, body.Avatar, body.Location, body.Coordinates, body.Fullname, body.Description, body.City, body.State, body.Country, body.ZipCode)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("could not update user")
	}

	return nil
}

func GetUserByUsername(context context.Context, username string) (Seller, error) {
	var userModel Seller
	validUserName := general.ValidUserName(username)
	if !validUserName {
		return userModel, errors.New("invalid username")
	}

	query := "SELECT * FROM sellers WHERE username = $1"
	err := database.GetContext(context, &userModel, query, username)
	if err == nil {
		return userModel, nil
	}
	return userModel, err
}

func CheckIfUsernameExists(context context.Context, username string) (Seller, error) {
	var user Seller
	user, err := GetUserByUsername(context, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New(messages.UserNotFoundMessage)
		}
		return user, err
	}

	return user, nil
}

func ResetOtpAndOtpExpiration(context context.Context, username string) error {
	_, err := database.ExecContext(context, "UPDATE sellers SET otp = '', otp_expiration = NULL WHERE username = $1", username)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSellerAccountStatus(context context.Context, username string, status string) error {
	_, err := database.ExecContext(context, "UPDATE sellers SET account_status = $2 WHERE username = $1", username, status)
	if err != nil {
		return err
	}
	return nil
}

func GetProductsListPaginatedValue(itemsPerPage, offset int, sid string) (*sql.Rows, error) {
	query := `SELECT id, title, views, price FROM products  WHERE seller_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	return database.Query(query, sid, itemsPerPage, offset)
}

func VerifySellerAccount(context context.Context, username string) error {
	_, err := database.ExecContext(context, "UPDATE sellers SET is_verified = $2 WHERE username = $1", username, true)
	if err != nil {
		return err
	}
	return nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsValidUser(context context.Context, username, password string) (Seller, error) {
	var userModel Seller
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

func UpdatePassword(context context.Context, username, newPassword string) error {
	password, err := hashPassword(newPassword)
	if err != nil {
		return errors.New(messages.SomethingWentWrongMessage)
	}

	query := "UPDATE sellers SET password = $2 WHERE username = $1"
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

func GetOtpFromDB(context context.Context, username string) (string, error) {
	var OTPSecret string
	var OTPExpiration time.Time

	query := "SELECT otp, otp_expiration FROM sellers WHERE username = $1"
	err := database.QueryRowContext(context, query, username).Scan(&OTPSecret, &OTPExpiration)
	if err != nil {
		return OTPSecret, err
	}

	if time.Now().After(OTPExpiration) {
		return "", errors.New(messages.OTPExpiredMessage)
	}

	return OTPSecret, nil
}

func SaveOTPAndExpirationInDB(context context.Context, username, otp string, expiration any) error {
	query := "UPDATE sellers SET otp = $2, otp_expiration = $3 WHERE username = $1"
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

func DeleteSellerByUsername(username string) error {
	query := "DELETE FROM sellers WHERE username = $1"
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptHashingCost)
	return string(bytes), err
}
