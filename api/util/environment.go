package util

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Returns the value of GO_ENV environment variable.
func GetGoEnv() string {
	return os.Getenv("GO_ENV")
}

// Returns either gin.ErrorTypePrivate or gin.ErrorTypePublic if GO_ENV is "dev".
func GetGinErrorType() gin.ErrorType {
	ginType := gin.ErrorTypePrivate
	if GetGoEnv() == "dev" {
		ginType = gin.ErrorTypePublic
	}
	return ginType
}