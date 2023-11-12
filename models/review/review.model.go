package reviewModel

import (
	"context"
	"olx-clone/infra/db"
)

var database = db.Mgr.DBConn

func AddReview(context context.Context, userId, productId, rating, comment string) error {
	query := `INSERT INTO reviews(user_id, product_id, rating, comment) VALUES($1, $2, $3, $4)`
	_, err := database.ExecContext(context, query, userId, productId, rating, comment)
	if err != nil {
		return err
	}
	return nil
}
