package productModel

type CreateProductBody struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Location    string `json:"location" validate:"required"`
	Coordinates string `json:"coordinates" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	SellerId    string `json:"seller_id" validate:"required"`
}

type ProductUpdateStruct struct {
	Title       string `validate:"required" json:"title"`
	Description string `validate:"required" json:"description"`
	Location    string `validate:"required" json:"location"`
	Coordinates string `validate:"required" json:"coordinates"`
	Price       string `validate:"required" json:"price"`
}

type ProductModel struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Location    string `json:"location" db:"location"`
	Coordinates string `json:"coordinates" db:"coordinates"`
	Price       string `json:"price" db:"price"`
	SellerId    string `json:"seller_id" db:"seller_id"`
	Views       string `json:"views" db:"views"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}
