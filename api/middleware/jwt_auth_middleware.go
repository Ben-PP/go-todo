package middleware

import (
	"fmt"
	"go-todo/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := util.DecodeTokenFromHeader(c)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "unauthorized",
				"detail": "Access token validation failed.",
			})
			c.Abort()
			return
		}

		c.Set("x-token-username", token.Username)
		c.Set("x-token-user-id", token.Subject)
		c.Set("x-token-is-admin", token.IsAdmin)

		c.Next()
	}
}