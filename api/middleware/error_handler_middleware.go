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

var ErrTokenValidationFailed = errors.New("token-validation-failed")

type ResponseParams struct {
	Status int
	StatusMessage string
	Detail string
}

type StatusMessages struct {
	Unauthorized		string
	InternalServerError	string
}
var statusMessages = StatusMessages{
	InternalServerError: "internal-server-error",
	Unauthorized: "unauthorized",
}

func getErrMeta(err *gin.Error) (*util.ErrorMeta) {
	var meta *util.ErrorMeta
	if errMeta, ok := err.Meta.(util.ErrorMeta); ok {
		meta = &errMeta
	}
	return meta
}

func metaPanic(meta *util.ErrorMeta, err *gin.Error) {
	if meta == nil {
		_, file, line, _ := runtime.Caller(1)
		logging.LogError(
			errors.New("err.Meta was not provided"),
			fmt.Sprintf("%v: %d", file, line),
			err.Error(),
		)
		panic("err.Meta is nil!")
	}
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		isPublic := err.Type == gin.ErrorTypePublic
		var params *ResponseParams
		meta := getErrMeta(err)
		switch {
		case errors.Is(err.Err, jwt.ErrTokenExpired):
			params = &ResponseParams{401,statusMessages.Unauthorized,err.Error()}
		case errors.Is(err.Err, jwt.ErrSignatureInvalid):
			params = &ResponseParams{401,statusMessages.Unauthorized, "invalid-token"}
		case errors.Is(err.Err, ErrTokenValidationFailed):
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{401, statusMessages.Unauthorized, "invalid-token"}
		default:
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{500, statusMessages.InternalServerError, ""}
			if isPublic {
				params.Detail = fmt.Sprintf("%v", meta.OrigErrMessage)
			}
		/*case err.Error() == "internal-server-error":
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{500, statusMessages.InternalServerError, ""}*/
		}
		
		body := gin.H{"status": params.StatusMessage}
		if params.Detail != "" {
			body = gin.H{
				"status": params.StatusMessage,
				"detail": params.Detail,
			}
		}
		c.JSON(params.Status, body)
	}
}