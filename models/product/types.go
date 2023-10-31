package productModel

type CreateProductBody struct {
	Title string `json:"title"`
	// Images      []string `json:"images"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Coordinates string `json:"coordinates"`
	Price       int    `json:"price"`
	SellerId    string `json:"seller_id"`
}

type ProductModel struct {
	Id          int      `json:"id"`
	Title       int      `json:"title"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Coordinates string   `json:"coordinates"`
	Price       int      `json:"price"`
	// SellerId    string      `json:"seller_id"`
	Views       int      `json:"views"`
	CreatedAt   string   `json:"created_at"`
}
