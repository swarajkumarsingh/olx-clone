package favoriteModel

type CreateUserStruct struct {
	UserId    string `validate:"required" json:"user_id"`
	ProductId string `validate:"required" json:"product_id"`
}
