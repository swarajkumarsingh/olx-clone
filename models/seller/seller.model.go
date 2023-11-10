package sellerModel

import (
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
