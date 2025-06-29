package middleware

import (
	"errors"
	"fmt"
	gterrors "go-todo/gt_errors"
	"go-todo/logging"
	"go-todo/util"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ResponseParams struct {
	Status int
	StatusMessage string
	Detail string
}

type StatusMessages struct {
	InternalServerError	string
	InvalidCredentials	string
	MalformedBody		string
	Unauthorized		string
}
var statusMessages = StatusMessages{
	InternalServerError: "internal-server-error",
	InvalidCredentials: "invalid-credentials",
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

// Handles the errors passed to the gin context and responds accordingly.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()
		isPublic := err.Type == gin.ErrorTypePublic
		var params *ResponseParams
		var pgErr *pgconn.PgError
		var authError *gterrors.GtAuthError
		var internalError *gterrors.GtInternalError
		switch {
		// Malformed requests
		case err.Type == gin.ErrorTypeBind:
			params = &ResponseParams{
				400,
				statusMessages.MalformedBody,
				err.Error(),
			}
		// GtAuthError
		case errors.As(err.Err, &authError):
			detail := ""
			status := 401
			var statusMessage string
			if isPublic {
				detail = authError.Reason.String()
			}
			switch authError.Reason {
			case gterrors.GtAuthErrorReasonExpired:
				statusMessage = statusMessages.Unauthorized
			case gterrors.GtAuthErrorReasonInvalidCredentials:
				statusMessage = statusMessages.InvalidCredentials
			case gterrors.GtAuthErrorReasonInvalidSignature:
				statusMessage = statusMessages.Unauthorized
			case gterrors.GtAuthErrorReasonTokenInvalid:
				statusMessage = statusMessages.Unauthorized
			case gterrors.GtAuthErrorReasonTokenReuse:
				statusMessage = statusMessages.Unauthorized
			case gterrors.GtAuthErrorReasonUsernameInvalid:
				statusMessage = statusMessages.InvalidCredentials
			case gterrors.GtAuthErrorReasonInternalError:
				statusMessage = statusMessages.InternalServerError
				detail = authError.Err.Error()
				logging.LogError(authError.Err, "unknown", authError.Error())	
			default:
				statusMessage = statusMessages.InternalServerError
				status = 500
				detail = authError.Err.Error()
			}

			params = &ResponseParams{
				Status: status,
				StatusMessage: statusMessage,
				Detail: detail,
			}
		// GtInternalError
		case errors.As(err, &internalError):
			logging.LogError(internalError.Err, internalError.File, "")
			detail := ""
			if isPublic {
				detail = internalError.Error()
			}
			var statusMessage string
			switch internalError.ResponseStatus {
			case 500:
				statusMessage = statusMessages.InternalServerError
			default:
				statusMessage = statusMessages.InternalServerError
			}
			params = &ResponseParams{
				Status: internalError.ResponseStatus,
				StatusMessage: statusMessage,
				Detail: detail,
			}
		// Jwt decoding errors
		/*case errors.Is(err.Err, jwt.ErrSignatureInvalid):
			logging.LogSecurityEvent(logging.SecurityScoreLow, logging.SecurityEventInvalidTokenSignature, "jwt-tokens")
			params = &ResponseParams{401,statusMessages.Unauthorized, "token-invalid"}
		case errors.Is(err.Err, jwt.ErrTokenMalformed):
			params = &ResponseParams{400, statusMessages.MalformedBody, "token-malformed"}
		case errors.Is(err.Err, ErrTokenValidationFailed):
			meta := getErrMeta[util.ErrInternalMeta](err)
			metaPanic(meta, err)
			logging.LogError(err.Err, meta.File, meta.OrigErrMessage)
			params = &ResponseParams{401, statusMessages.Unauthorized, "token-invalid"}*/
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