package authentication

import (
	"net/http"
	"olx-clone/conf"
	userModel "olx-clone/models/user"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorize(ctx *gin.Context) {
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

	// check if the user exists
	// check user role(user, seller, admin)
	// if seller account, check if the account is not suspended

	ctx.Set("username", claims.Username)

	ctx.Next()
}
