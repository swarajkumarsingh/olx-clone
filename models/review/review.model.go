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

func GetReview(context context.Context, rid string) (Review, error) {
	var review Review
	query := "SELECT * FROM reviews WHERE id = $1"
	err := database.GetContext(context, &review, query, rid)
	if err == nil {
		return review, nil
	}
	return review, err
}
