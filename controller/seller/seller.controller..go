package seller

import (
	"context"
	"database/sql"
	"net/http"
	"olx-clone/constants"
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
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username := ctx.Param("sid")
	if username == "" {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidUsernameMessage)
	}

	seller, err := model.GetSellerByUsername(context.TODO(), username)
	if err == sql.ErrNoRows {
		logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.SellerNotFoundMessage)
	}
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"user":  seller,
	})
}

// login
func LoginSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.LoginUser
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(messages.InvalidUsernameOrPasswordMessage)
	}
	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(messages.InvalidUsernameOrPasswordMessage)
	}

	_, err := model.IsValidUser(context.TODO(), body.Username, body.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.UserNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(err)
	}

	token, err := generateJwtToken(body.Username)
	if err != nil {
		logger.WithRequest(ctx).Panicln("unable to login, try again later")
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Seller login successfully",
		"token":   token,
	})
}

// update
func UpdateSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.SellerUpdateStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(messages.InvalidBodyMessage)
	}

	username := ctx.Param("sid")

	if err := model.UpdateSeller(context.TODO(), username, body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Seller updated successfully",
	})
}

// delete
func DeleteSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username := ctx.Param("sid")
	if username == "" {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidUsernameMessage)
	}

	if err := model.DeleteSellerByUsername(username); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

// request reset password
func RequestResetPassword(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.ResetRequestStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(messages.InvalidBodyMessage)
	}

	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, err)
	}

	user, err := model.CheckIfUsernameExists(context.TODO(), body.Username)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	otp, err := generateOtp()
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	_, err = sendOtpEmail(constants.DefaultSenderEmailId, user.Email, otp)
	if err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	otpSecret := encodeString(otp)
	otpExpiration := getTimeInMinutes(5)
	if err = model.SaveOTPAndExpirationInDB(context.TODO(), user.Username, otpSecret, otpExpiration); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "OTP send successfully",
	})
}

// reset password
func ResetPasswordSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	var body model.ResetPasswordStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.WithRequest(ctx).Panicln(messages.InvalidBodyMessage)
	}

	if err := general.ValidateStruct(body); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	encodedOTP, err := model.GetOtpFromDB(context.TODO(), body.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.WithRequest(ctx).Panicln(http.StatusNotFound, messages.UserNotFoundMessage)
		}
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	otp, err := decodeString(encodedOTP)
	if err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	valid := verifyOTPs(body.OTP, otp)
	if !valid {
		logger.WithRequest(ctx).Panicln(http.StatusBadRequest, messages.InvalidOTPMessage)
	}

	if err := model.UpdatePassword(context.TODO(), body.Username, body.NewPassword); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	if err := model.ResetOtpAndOtpExpiration(context.TODO(), body.Username); err != nil {
		logger.WithRequest(ctx).Panicln(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Reset password successful",
	})
}

// suspend seller account
func SuspendSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username := ctx.Param("username")
	if err := model.UpdateSellerAccountStatus(context.TODO(), username, constants.StatusSuspendSeller); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "seller suspended successfully",
	})
}

func ReportSeller(ctx *gin.Context) {

}

// ban seller account
func BanSeller(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username := ctx.Param("username")
	if err := model.UpdateSellerAccountStatus(context.TODO(), username, constants.StatusBanSeller); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "seller suspended successfully",
	})
}

// ban seller account
func ActivateSellerAccount(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

	username := ctx.Param("username")
	if err := model.UpdateSellerAccountStatus(context.TODO(), username, constants.StatusActiveSeller); err != nil {
		logger.WithRequest(ctx).Panicln(http.StatusInternalServerError, messages.SomethingWentWrongMessage)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "seller suspended successfully",
	})
}

// get all created products
func GetAllCreatedProduct(ctx *gin.Context) {

}

// verify seller account
func VerifySeller(ctx *gin.Context) {

}
