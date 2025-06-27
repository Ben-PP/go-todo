package util

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyCustomClaims struct {
	IsAdmin bool `json:"is_admin"`
	Username string `json:"username"`
	Family string `json:"family"`
	jwt.RegisteredClaims
}

var globalConfig *Config

func getConfig() (*Config, error) {
	if globalConfig == nil {
		var err error
		config, err := LoadConfig(".")
		if err != nil {
			return nil, err
		}
		globalConfig = &config
	}

	return globalConfig, nil
}

func generateToken(username string, userID string, isAdmin bool, isRefreshToken bool, family string) (string, *MyCustomClaims, error) {
	config, err := getConfig()
	if err != nil {
		return "", nil, err
	}
	
	authLifeSpan := config.AccessTokenLifeSpan
	if isRefreshToken {
		authLifeSpan = config.RefreshTokenLifeSpan
	}
	lifeSpanDuration := time.Minute * time.Duration(authLifeSpan)
	timeNow := time.Now().UTC()
	expiry := jwt.NewNumericDate(timeNow.Add(lifeSpanDuration))	
	if family == "" && isRefreshToken {
		family = uuid.New().String()
	} else if !isRefreshToken {
		family = "access"
	}
	claims := MyCustomClaims{
		isAdmin,
		username,
		family,
		jwt.RegisteredClaims{
			ID: uuid.New().String(),
			Subject: userID,
			ExpiresAt: expiry,
			Issuer: "GO-TODO",
			IssuedAt: jwt.NewNumericDate(timeNow),

		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	secret := config.JwtAccessSecret
	if isRefreshToken {
		secret = config.JwtRefreshSecret
	}
	encodedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}
	return encodedToken, &claims, nil
}

func GenerateAccessToken(username string, userID string, isAdmin bool) (string, *MyCustomClaims, error) {
	return generateToken(username, userID, isAdmin, false, "")
}

func GenerateRefreshToken(username string, userID string, isAdmin bool, tokenFamily string) (string, *MyCustomClaims, error) {
	return generateToken(username, userID, isAdmin, true, tokenFamily)
}

func decodeToken(tokenString string, isRefreshToken bool) (*MyCustomClaims, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	secret := config.JwtAccessSecret
	if isRefreshToken {
		secret = config.JwtRefreshSecret
	}
	decodedToken, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token)(any, error){
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := decodedToken.Claims.(*MyCustomClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("something went wrong decoding token")
	}
}

func DecodeAccessToken(tokenString string) (*MyCustomClaims, error) {
	return decodeToken(tokenString, false)
}

func DecodeRefreshToken(tokenString string) (*MyCustomClaims, error) {
	return decodeToken(tokenString, true)
}

func getTokenFromHeader(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func DecodeTokenFromHeader(c *gin.Context) (*MyCustomClaims, error) {
	tokenString := getTokenFromHeader(c)
	return DecodeAccessToken(tokenString)
}