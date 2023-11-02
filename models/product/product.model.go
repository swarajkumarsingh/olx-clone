package productModel

import (
	"context"
	"errors"
	"olx-clone/constants/messages"
	"olx-clone/infra/db"
)

var database = db.Mgr.DBConn

func CreateProduct(context context.Context, body CreateProductBody) error {
	query := "INSERT INTO  products(title, description, location, coordinates, price, seller_id) VALUES($1, $2, $3, $4, $5, $6)"
	res, err := database.ExecContext(context, query, body.Title, body.Description, body.Location, body.Coordinates, body.Price, body.SellerId)
	if err != nil {
		return errors.New(messages.SomethingWentWrongMessage)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 || err != nil {
		return errors.New(messages.SomethingWentWrongMessage)
	}

	return nil
}

func GetProduct(context context.Context, productId string) (ProductModel, error) {
	var productModel ProductModel
	query := "SELECT * FROM products WHERE id = $1"
	err := database.GetContext(context, &productModel, query, productId)
	if err != nil {
		return productModel, err
	}
	return productModel, nil
}

func AddProductViews(context context.Context, userId, productId string) error {
	query := "INSERT INTO  product_views(user_id, product_id) VALUES($1, $2)"
	res, err := database.ExecContext(context, query, userId, productId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("unable to add product views")
	}
	return nil
}
