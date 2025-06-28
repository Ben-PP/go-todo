package middleware

import (
	"errors"
	"fmt"
	"go-todo/logging"
	"go-todo/util"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := util.DecodeTokenFromHeader(c)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				fmt.Println("HEREEEEEEEEEEEEEEEE")
			}
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrSignatureInvalid) {
				logging.LogTokenUsage(false, "use", "access", token)

				var errType gin.ErrorType
				if errors.Is(err, jwt.ErrTokenExpired) {
					errType = gin.ErrorTypePublic
				} else {
					errType = gin.ErrorTypePrivate
				}

				c.Error(err).SetType(errType)
				c.Abort()
				return
			} else {
				logging.LogTokenUsage(
					false,
					"use",
					"access",
					token,
				)
				_, file, line, _ := runtime.Caller(1)
				c.Error(errors.New("token-validation-error")).
				SetType(gin.ErrorTypePrivate).SetMeta(util.ErrorMeta{
					File: fmt.Sprintf("%v: %d", file, line),
					OrigErrMessage: err.Error(),
				})
			}
			c.Abort()
			return
		}

		logging.LogTokenUsage(true, "use", "access", token)

		c.Set("x-token-username", token.Username)
		c.Set("x-token-user-id", token.Subject)
		c.Set("x-token-is-admin", token.IsAdmin)

		c.Next()
	}
}