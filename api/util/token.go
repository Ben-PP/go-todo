package util

import (
	db "go-todo/db/sqlc"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateToken(user db.User) (string, error) {
	config, err := LoadConfig(".")
	if err != nil {
		return "", err
	}
	// TODO Continue here. Create refresh logic
	authLifeSpan := config.AuthTokenLifeSpan
	/*refreshLifeSpan, err := util.Config.RefreshTokenLifeSpan
	if err != nil {
		return "", err
	}*/

	claims := MyCustomClaims{
		user.IsAdmin,
		jwt.RegisteredClaims{
			ID: user.ID,
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(
					time.Minute * time.Duration(authLifeSpan),
				),
			),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(config.JwtSecret))
}