package middleware

import (
	"errors"
	"fmt"
	"go-todo/logging"
	"go-todo/util"
	"runtime"

	"github.com/gin-gonic/gin"
)

type ResponseParams struct {
	Status int
	StatusMessage string
	Detail string
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

		switch err.Error() {
		case "token-expired":
			params = &ResponseParams{401,"unauthorized",err.Error()}
		case "token-signature-invalid":
			params = &ResponseParams{401,"unauthorized", "invalid-token"}
		case "token-validation-error":
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{401, "unauthorized", "invalid-token"}
		case "internal-server-error":
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{500, err.Error(), ""}
			if isPublic {
				params.Detail = fmt.Sprintf("%v", err.Meta)
			}
		default:
			metaPanic(meta, err)
			logging.LogError(err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{500, "internal-server-error", ""}
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