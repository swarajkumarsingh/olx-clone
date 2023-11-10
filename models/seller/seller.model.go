package sellerModel

import (
	"context"
	"database/sql"
	"errors"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	"olx-clone/infra/db"
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