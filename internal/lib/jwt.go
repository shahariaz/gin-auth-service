package lib

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint, username, role string, secret []byte) (string, error) {
	claims := TokenClaims{
		Username: username,
		Role:     role,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-gin-app",
			Audience:  []string{"api"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(userID uint, username string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iss":      "my-gin-app",
		"aud":      "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
