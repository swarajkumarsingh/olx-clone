package favoriteModel

import (
	"context"
	"database/sql"
	"errors"
	"olx-clone/constants/messages"
	"olx-clone/infra/db"
)

var database = db.Mgr.DBConn

func AddFavorite(context context.Context, userId, productId string) error {
	query := `INSERT INTO favorites(user_id, product_id) VALUES($1, $2)`
	result, err := database.ExecContext(context, query, userId, productId)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return errors.New(messages.SomethingWentWrongMessage)
	}
	return nil
}

func RemoveFavorite(context context.Context, fid string) error {
	query := "DELETE FROM favorites WHERE id = $1"
	res, err := database.ExecContext(context, query, fid)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("favorite already deleted or favorite not found")
	}

	return err
}

func GetUsersListPaginatedValue(itemsPerPage, offset int) (*sql.Rows, error) {
	query := `SELECT id, user_id, product_id FROM favorites ORDER BY id LIMIT $1 OFFSET $2`
	return database.Query(query, itemsPerPage, offset)
}
