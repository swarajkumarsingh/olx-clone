package productModel

import (
	"context"
	"database/sql"
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

func UpdateProduct(context context.Context, productId string, body ProductUpdateStruct) error {
	query := "UPDATE products SET title = $2, description = $3, location = $4, coordinates = $5, price = $6 WHERE id = $1"
	res, err := database.ExecContext(context, query, productId, body.Title, body.Description, body.Location, body.Coordinates, body.Price)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("could not update user")
	}

	return nil
}

func DeleteProductByProductId(context context.Context, productId string) error {
	query := "DELETE FROM products WHERE id = $1"
	res, err := database.ExecContext(context, query, productId)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product already deleted or product not found")
	}

	return err
}

func DeleteProductViewsByProductId(context context.Context, productId string) error {
	query := "DELETE FROM product_views WHERE product_id = $1"
	_, err := database.ExecContext(context, query, productId)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func GetUsersListPaginatedValue(itemsPerPage, offset int) (*sql.Rows, error) {
	query := `SELECT id, title, views, price FROM products ORDER BY id LIMIT $1 OFFSET $2`
	return database.Query(query, itemsPerPage, offset)
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
