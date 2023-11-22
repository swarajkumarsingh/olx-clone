package authentication

import (
	"context"
	"database/sql"
	"net/http"
	"olx-clone/conf"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		ctx.Abort()
		return
	}

	// Token format: Bearer <token>
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		ctx.Abort()
		return
	}

	tokenString := splitToken[1]

	token, err := jwt.ParseWithClaims(tokenString, &userModel.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.JWTSecretKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(*userModel.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		ctx.Abort()
		return
	}

	username := claims.Username

	// check if the user exists
	_, err = userModel.CheckIfUsernameExists(context.TODO(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		ctx.Abort()
		return
	}

	ctx.Set("username", claims.Username)

	ctx.Next()
}

func AuthorizeSeller(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		ctx.Abort()
		return
	}

	// Token format: Bearer <token>
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		ctx.Abort()
		return
	}

	tokenString := splitToken[1]

	token, err := jwt.ParseWithClaims(tokenString, &userModel.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.JWTSecretKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
		ctx.Abort()
		return
	}

	claims, ok := token.Claims.(*userModel.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		ctx.Abort()
		return
	}

	username := claims.Username

	// check if the user exists
	_, err = sellerModel.CheckIfUsernameExists(context.TODO(), username)
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

	ctx.Set("username", claims.Username)
	ctx.Next()
}
