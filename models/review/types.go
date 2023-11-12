package reviewModel

type CreateReviewStruct struct {
	Rating    string `validate:"required" json:"rating"`
	UserId    string `validate:"required" json:"user_id"`
	Comment   string `validate:"required" json:"comment"`
	ProductId string `validate:"required" json:"product_id"`
}
type Review struct {
	Id        int    `json:"id" db:"id"`
	UserId    string `json:"user_id" db:"user_id"`
	ProductId string `json:"product_id" db:"product_id"`
	Rating    string `json:"rating" db:"rating"`
	Comment   string `json:"comment" db:"comment"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type UpdateReviewStruct struct {
	Rating   string `validate:"required" json:"rating"`
	Comment  string `validate:"required" json:"comment"`
	ReviewId string `validate:"required" json:"review_id"`
}
