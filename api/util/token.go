package util

import (
	db "go-todo/db/sqlc"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyCustomClaims struct {
	IsAdmin bool `json:"is_admin"`
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func generateToken(user db.User, isRefreshToken bool) (string, error) {
	config, err := LoadConfig(".")
	if err != nil {
		return "", err
	}
	
	authLifeSpan := config.AccessTokenLifeSpan
	if isRefreshToken {
		authLifeSpan = config.RefreshTokenLifeSpan
	}

	claims := MyCustomClaims{
		user.IsAdmin,
		user.Username,
		jwt.RegisteredClaims{
			ID: uuid.New().String(),
			Subject: user.ID,
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(
					time.Minute * time.Duration(authLifeSpan),
				),
			),
			Issuer: "GO-TODO",

		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(config.JwtSecret))
}

func GenerateAccessToken(user db.User) (string, error) {
	return generateToken(user, false)
}

func GenerateRefreshToken(user db.User) (string, error) {
	return generateToken(user, true)
}