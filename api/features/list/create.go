package list

import (
	"errors"
	"fmt"
	db "go-todo/db/sqlc"
	"go-todo/gterrors"
	"go-todo/logging"
	"go-todo/schemas"
	"go-todo/util/mycontext"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (controller *ListController) CreateList(ctx *gin.Context) {
	var payload *schemas.CreateList
	if ok := mycontext.ShouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	}

	tokenUserId, tokenUserName, _, err := mycontext.GetTokenVariables(ctx)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		mycontext.CtxAddGtInternalError("failed to get claims from jwt", file, line, err, ctx)
		return
	}

	args := &db.CreateListParams{
		ID: uuid.New().String(),
		UserID: tokenUserId,
		Title: payload.Title,
		Description: payload.Description,
	}

	list, err := controller.db.CreateList(ctx, *args)
	if err != nil {
		var pgErr *pgconn.PgError
		errMessage := "failed to create list"
		if errors.As(err, &pgErr) {
			_, file, line, _ := runtime.Caller(0)
			switch pgErr.Code {
			case "23505":
				mycontext.CtxAddGtInternalError(
					"failed to create unique id for list",
					file,
					line,
					err,
					ctx,
				)
			default:
				_, file, line, _ := runtime.Caller(0)
				mycontext.CtxAddGtInternalError(errMessage, file, line, err, ctx)
			}
			return
		}
		_, file, line, _ := runtime.Caller(0)
		mycontext.CtxAddGtInternalError(errMessage, file, line, err, ctx)
		return
	}

	reqUser, err := controller.db.GetUserById(ctx, tokenUserId)
	if err != nil {
		logging.LogSecurityEvent(
			logging.SecurityScoreLow,
			logging.SecurityEventJwtUserUnknown,
			ctx.FullPath(),
			tokenUserName,
			ctx.ClientIP(),
		)
		ctx.Error(
			gterrors.NewGtAuthError(
				gterrors.GtAuthErrorReasonJwtUserNotFound,
				fmt.Errorf("could not get user from db: %w", err),
			),
		).SetType(gterrors.GetGinErrorType())
		return
	}

	logging.LogObjectEvent(
		ctx.FullPath(),
		ctx.ClientIP(),
		logging.ObjectEventCreate,
		&reqUser,
		&list,
		nil,
		logging.ObjectEventSubList,
	)
	ctx.JSON(200, gin.H{"status": "created", "list": list})
}