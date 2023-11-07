package favoriteModel

import (
	"context"
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
