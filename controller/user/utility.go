package user

import (
	"errors"
	"olx-clone/infra/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("chotu_babu_is_not_chotu_any_more")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CalculateTotalPages(page, itemsPerPage int) int {
	if page <= 0 {
		return 1
	}
	return (page + itemsPerPage - 1) / itemsPerPage
}

func InsertUser(ctx *gin.Context, body UserBody, hashedPassword string) error {
	database := db.Mgr.DBConn
	query := `INSERT INTO "user"(name, email, password, number) VALUES($1, $2, $3, $4)`
	_, err := database.Exec(query, body.Name, body.Email, hashedPassword, body.Number)
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJwtToken(name string) (string, error) {
	expirationTime := time.Now().Add(5 * 24 * time.Hour)
	claims := &Claims{
		Username: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	if tokenString == "" {
		return "", errors.New("error while authorizing")
	}

	return tokenString, nil
}
