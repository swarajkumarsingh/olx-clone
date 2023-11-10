package seller

import (
	"net/http"
	"olx-clone/constants/messages"
	"olx-clone/errorHandler"
	"olx-clone/functions/general"
	"olx-clone/functions/logger"
	model "olx-clone/models/seller"

	"github.com/gin-gonic/gin"
)

// create seller
func CreateSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.SellerBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidBodyMessage)
	}

	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	if model.SellerAlreadyExistsWithUsername(body.Username) {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.SellerAlreadyExistsMessage)
	}

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	if err = model.CreateSeller(body, hashedPassword); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := generateJwtToken(body.Username)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Seller Created successfully",
		"token":   token,
	})
}

// get all seller - admin
func GetAllSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	page := getCurrentPageValue(ctx)
	itemsPerPage := getItemPerPageValue(ctx)
	offset := getOffsetValue(page, itemsPerPage)

	rows, err := model.GetSellerListPaginatedValue(itemsPerPage, offset)
	if err != nil {
		logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
	}
	defer rows.Close()

	sellers := make([]gin.H, 0)

	for rows.Next() {
		var id int
		var username, email, number string
		if err := rows.Scan(&id, &username, &email, &number); err != nil {
			logger.WithRequest(ctx).Panicln(messages.FailedToRetrieveUsersMessage)
		}
		sellers = append(sellers, gin.H{"id": id, "username": username, "email": email, "number": number})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":       false,
		"sellers":     sellers,
		"page":        page,
		"per_page":    itemsPerPage,
		"total_pages": calculateTotalPages(page, itemsPerPage),
	})
}

// get user
func GetSeller(ctx *gin.Context) {

}

// login
func LoginSeller(ctx *gin.Context) {

}

// logout
func LogoutSeller(ctx *gin.Context) {

}

// update
func UpdateSeller(ctx *gin.Context) {

}

// delete
func DeleteSeller(ctx *gin.Context) {

}

// request reset password
func RequestResetPasswordSeller(ctx *gin.Context) {

}

// reset password
func ResetPasswordSeller(ctx *gin.Context) {

}

// suspend seller account
func SuspendSeller(ctx *gin.Context) {

}

// un-suspend seller account
func UnSuspendSeller(ctx *gin.Context) {

}

// ban seller account
func BanSeller(ctx *gin.Context) {

}

// ban seller account
func UnBanSeller(ctx *gin.Context) {

}

// get all created products
func GetAllCreatedProduct(ctx *gin.Context) {

}

// verify seller account
func VerifySeller(ctx *gin.Context) {

}
