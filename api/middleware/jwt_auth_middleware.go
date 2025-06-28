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
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrSignatureInvalid) {
				logging.LogTokenEvent(false, c.FullPath(), logging.TokenEventTypeUse, c.RemoteIP(), token)

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
				logging.LogTokenEvent(
					false,
					c.FullPath(),
					logging.TokenEventTypeUse,
					c.RemoteIP(),
					token,
				)
				_, file, line, _ := runtime.Caller(1)
				c.Error(ErrTokenValidationFailed).
				SetType(gin.ErrorTypePrivate).SetMeta(util.ErrInternalMeta{
					File: fmt.Sprintf("%v: %d", file, line),
					OrigErrMessage: err.Error(),
				})
			}
			c.Abort()
			return
		}

		logging.LogTokenEvent(true, c.FullPath(), logging.TokenEventTypeUse, c.RemoteIP(), token)

		c.Set("x-token-username", token.Username)
		c.Set("x-token-user-id", token.Subject)
		c.Set("x-token-is-admin", token.IsAdmin)

		c.Next()
	}
}