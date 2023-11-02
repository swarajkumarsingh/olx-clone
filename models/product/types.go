package productModel

type CreateProductBody struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Coordinates string `json:"coordinates"`
	Price       int    `json:"price"`
	SellerId    string `json:"seller_id"`
}

type ProductModel struct {
	Id          int      `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	Description string   `json:"description" db:"description"`
	Location    string   `json:"location" db:"location"`
	Coordinates string   `json:"coordinates" db:"coordinates"`
	Price       string   `json:"price" db:"price"`
	SellerId    string   `json:"seller_id" db:"seller_id"`
	Views       string   `json:"views" db:"views"`
	CreatedAt   string   `json:"created_at" db:"created_at"`
}
