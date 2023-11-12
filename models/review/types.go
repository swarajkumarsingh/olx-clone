package reviewModel

type CreateReviewStruct struct {
	Rating    string `validate:"required" json:"rating"`
	UserId    string `validate:"required" json:"user_id"`
	Comment   string `validate:"required" json:"comment"`
	ProductId string `validate:"required" json:"product_id"`
}
