package authentication

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"olx-clone/conf"
	"olx-clone/constants"
	"olx-clone/constants/messages"
	sellerModel "olx-clone/models/seller"
	userModel "olx-clone/models/user"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthorizeUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Authorization header is missing"})
		ctx.Abort()
		return
	}

	// Token format: Bearer <token>
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token format"})
		ctx.Abort()
		return
	}

	tokenString := splitToken[1]

	token, err := jwt.ParseWithClaims(tokenString, &userModel.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.JWTSecretKey, nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token"})
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(*userModel.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token claims"})
		ctx.Abort()
		return
	}

	userId := claims.UserId

	// check if the user exists
	_, err = userModel.CheckIfUsernameExistsWithId(context.TODO(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token claims 2"})
			ctx.Abort()
			return
		}
		fmt.Println(userId, claims.ID, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token claims 3"})
		ctx.Abort()
		return
	}

	ctx.Set(constants.UserIdMiddlewareConstant, userId)
	ctx.Next()
}

func AuthorizeSeller(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Authorization header is missing"})
		ctx.Abort()
		return
	}

	// Token format: Bearer <token>
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token format"})
		ctx.Abort()
		return
	}

	tokenString := splitToken[1]

	token, err := jwt.ParseWithClaims(tokenString, &userModel.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.JWTSecretKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token"})
		ctx.Abort()
		return
	}

	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token"})
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(*userModel.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Invalid token claims"})
		ctx.Abort()
		return
	}

	userId := claims.UserId

	// check if the user exists
	_, err = sellerModel.CheckIfUsernameExistsWithId(context.TODO(), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": messages.SellerNotFoundMessage})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": messages.SomethingWentWrongMessage})
		ctx.Abort()
		return
	}

	ctx.Set(constants.SellerIdMiddlewareConstant, claims.UserId)
	ctx.Next()
}
