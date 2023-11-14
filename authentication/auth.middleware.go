package authentication

import (
	"net/http"
	"olx-clone/conf"
	"olx-clone/errorHandler"
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

	// Check if the token is valid
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
		ctx.Abort()
		return
	}

	// Access the username from the token claims
	claims, ok := token.Claims.(*userModel.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		ctx.Abort()
		return
	}

	// Save the username in the request context
	ctx.Set("username", claims.Username)

	// Continue processing the request
	ctx.Next()
}

func TokenAuthenticate(ctx *gin.Context) {
	defer errorHandler.Recovery(ctx, http.StatusConflict)

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

	// Parse and verify the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return conf.JWTSecretKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}

	// Check if the token is valid
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
		ctx.Abort()
		return
	}

	// Extract user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		ctx.Abort()
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		ctx.Abort()
		return
	}

	// Save the user ID in the request context
	ctx.Set("userId", int(userID))

	// Continue processing the request
	ctx.Next()
}
