package controllers

import (
	"fmt"
	gterrors "go-todo/gt_errors"
	"go-todo/util"

	"github.com/gin-gonic/gin"
)

// Check the body format against the schema. 'payload' should be like
// &schemas.<MyBodySchema>
func shouldBindBodyWithJSON(payload any, c *gin.Context) bool {
	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return false
	}
	return true
}

func ctxAddGtInternalError(message, file string, line int, err error, c *gin.Context) {
	errToAdd := err
	if message != "" {
		errToAdd = fmt.Errorf("%v: %w", message, err)
	}
	c.Error(
		gterrors.NewGtInternalError(
			errToAdd,
			util.GetFileNameWithLine(file,line),
			500,
		),
	).SetType(util.GetGinErrorType())
}