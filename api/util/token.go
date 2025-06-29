package util

import (
	"errors"
	"fmt"
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

type JwtValidationError struct  {
	Claims *MyCustomClaims
	OrigErr error
	Msg string
}

func (e *JwtValidationError) Error() string {
	return fmt.Sprintf("failed to validate token: %v", e.OrigErr.Error())
}

func generateToken(username string, userID string, isAdmin bool, isRefreshToken bool, family string) (string, *MyCustomClaims, error) {
	config, err := GetConfig()
	if err != nil {
		return "", nil, err
	}
	
	authLifeSpan := 1//config.AccessTokenLifeSpan
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

// Takes a jwt as a string and boolean isRefreshToken telling should it be
// decoded with refresh secret. If all goes well, returns claims and if not,
// returns JwtValidationError or normal error.
func decodeToken(tokenString string, isRefreshToken bool) (*MyCustomClaims, error) {
	config, err := GetConfig()
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
		if claims, ok := decodedToken.Claims.(*MyCustomClaims); ok {
			return nil, &JwtValidationError{
				Claims: claims,
				OrigErr: err,
				Msg: fmt.Sprintf("error parsing jwt: %v", err.Error()),
			}
		}
		return nil, &JwtValidationError{
			Claims: &MyCustomClaims{},
			OrigErr: err,
			Msg: fmt.Sprintf("error parsing jwt: %v", err.Error()),
		}
	} else if claims, ok := decodedToken.Claims.(*MyCustomClaims); ok {
		return claims, nil
	} else {
		return nil, &JwtValidationError{
			Claims: &MyCustomClaims{},
			OrigErr: errors.New("something went wrong decoding token"),
			Msg: "error parsing jwt",
		}
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