package util

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func getBooleanKey(key string, c *gin.Context) (bool, bool) {
	var isAdmin bool
	isAdminRaw, exists := c.Get(key)
	if exists {
		if val, ok := isAdminRaw.(bool); ok {
			isAdmin = val
		} else {
			exists = false
		}
	}
	return isAdmin, exists
}

func GetTokenVariables(ctx *gin.Context) (userID string, username string, isAdmin bool, err error) {
	userID = ctx.GetString("x-token-user-id")
	username = ctx.GetString("x-token-username")
	isAdmin, isAdminExists := getBooleanKey("x-token-is-admin", ctx)

	if userID == "" || username == "" || !isAdminExists {
		return "", "", false, errors.New("failed to get token variables")
	}

	return userID, username, isAdmin, nil
}