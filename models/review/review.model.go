package reviewModel

import (
	"context"
	"database/sql"
	"errors"
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

func UpdateReview(context context.Context, reviewId, rating, comment string) error {
	query := "UPDATE reviews SET rating = $2, comment = $3 WHERE id = $1"
	res, err := database.ExecContext(context, query, reviewId, rating, comment)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("could not update user")
	}

	return nil
}

func DeleteReview(context context.Context, rid string) error {
	query := "DELETE FROM reviews WHERE id = $1"
	res, err := database.ExecContext(context, query, rid)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("review already deleted or not found")
	}

	return err
}

func GetProductReviews(context context.Context, pid string, itemsPerPage, offset int) (*sql.Rows, error) {
	query := `SELECT id, product_id, rating, comment FROM reviews WHERE product_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	return database.QueryContext(context, query, pid, itemsPerPage, offset)
}
