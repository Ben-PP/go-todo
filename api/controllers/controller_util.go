package controllers

import "github.com/gin-gonic/gin"

func shouldBindBodyWithJSON(payload any, c *gin.Context) bool {
	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return false
	}
	return true
}