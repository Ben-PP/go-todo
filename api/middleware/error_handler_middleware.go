package middleware

import (
	"errors"
	"fmt"
	"go-todo/logging"
	"go-todo/util"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrTokenValidationFailed = errors.New("token-validation-failed")

type ResponseParams struct {
	Status int
	StatusMessage string
	Detail string
}

type StatusMessages struct {
	InternalServerError	string
	MalformedBody		string
	Unauthorized		string
}
var statusMessages = StatusMessages{
	InternalServerError: "internal-server-error",
	MalformedBody: "malformed-body",
	Unauthorized: "unauthorized",
}

func getErrMeta[T any](err *gin.Error) (*T) {
	var meta *T
	if errMeta, ok := err.Meta.(T); ok {
		meta = &errMeta
	}
	return meta
}

func metaPanic[T any](meta *T, err *gin.Error) {
	if meta == nil {
		_, file, line, _ := runtime.Caller(1)
		logging.LogError(
			errors.New("MetaPanic: err.Meta was not provided"),
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
		isPublic := err.Type == gin.ErrorTypePublic || err.Type == gin.ErrorTypeBind
		var params *ResponseParams
		var pgErr *pgconn.PgError
		switch {
		// Malformed requests
		case err.Type == gin.ErrorTypeBind:
			params = &ResponseParams{
				400,
				statusMessages.MalformedBody,
				err.Error(),
			}
		// Jwt decoding errors
		case errors.Is(err.Err, jwt.ErrTokenExpired):
			params = &ResponseParams{401,statusMessages.Unauthorized,err.Error()}
		case errors.Is(err.Err, jwt.ErrSignatureInvalid):
			logging.LogSecurityEvent(logging.SecurityScoreLow, logging.SecurityEventInvalidTokenSignature, "jwt-tokens")
			params = &ResponseParams{401,statusMessages.Unauthorized, "token-invalid"}
		case errors.Is(err.Err, jwt.ErrTokenMalformed):
			params = &ResponseParams{400, statusMessages.MalformedBody, "token-malformed"}
		case errors.Is(err.Err, ErrTokenValidationFailed):
			meta := getErrMeta[util.ErrInternalMeta](err)
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{401, statusMessages.Unauthorized, "token-invalid"}
		// Jwt generating errors
		// Database errors
		case errors.Is(err.Err, pgx.ErrNoRows):
			meta := getErrMeta[util.ErrDatabaseMeta](err)
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.QueryDetails)
			var body string
			switch meta.RespondWithStatus {
			case 400:
				body = statusMessages.MalformedBody
			case 401:
				body = statusMessages.Unauthorized
			case 500:
				body = statusMessages.InternalServerError
			default:
				body = ""
			}
			params = &ResponseParams{meta.RespondWithStatus, body, ""}
		case errors.As(err.Err, &pgErr):
			meta := getErrMeta[util.ErrDatabaseMeta](err)
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.QueryDetails)
			var body string
			switch meta.RespondWithStatus {
			case 500:
				body = statusMessages.InternalServerError
			default:
				body = statusMessages.InternalServerError
			}
			params = &ResponseParams{500, body, ""}
		// Catch all
		default:
			meta := getErrMeta[util.ErrInternalMeta](err)
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{500, statusMessages.InternalServerError, ""}
			if isPublic {
				params.Detail = fmt.Sprintf("%v", meta.OrigErrMessage)
			}
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