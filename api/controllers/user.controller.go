package controllers

import (
	"context"
	"errors"
	"fmt"
	db "go-todo/db/sqlc"
	gterrors "go-todo/gt_errors"
	"go-todo/logging"
	"go-todo/schemas"
	"go-todo/util"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserController struct {
	db *db.Queries
	ctx context.Context
}

func NewUserController(db *db.Queries, ctx context.Context) *UserController {
	return &UserController{db:db, ctx: ctx}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var payload *schemas.CreateUser
	if ok := shouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	}

	makeAdmin := false
	users, err := uc.db.GetAllUsers(ctx)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("failed to get users from db", file, line, err, ctx)
		return
	}
	if len(users) == 0 {
		makeAdmin = true
	}

	isPasswdValid, err := util.ValidatePassword(payload.Password)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("failed to validate new password", file, line, err, ctx)
		return
	} else if !isPasswdValid {
		ctx.Error(gterrors.ErrPasswordUnsatisfied).SetType(gin.ErrorTypePublic)
		return
	}

	isUsernameValid, err := util.ValidateUsername(payload.Username)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("failed to validate new username", file, line, err, ctx)
		return
	} else if !isUsernameValid {
		ctx.Error(gterrors.ErrUsernameUnsatisfied).SetType(gin.ErrorTypePublic)
		return
	}

	userUUID := uuid.New()
	passwd := payload.Password
	passwdHash,err := util.HashPassword(passwd)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("failed to hash new password", file, line, err, ctx)
		return
	}

	args := &db.CreateUserParams{
		ID: userUUID.String(),
		Username: payload.Username,
		PasswordHash: passwdHash,
		IsAdmin: makeAdmin,
	}

	user, err := uc.db.CreateUser(ctx, *args)
	if err != nil {
		var pgErr *pgconn.PgError
		errMessage := "failed to create user"
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				ctx.Error(gterrors.ErrUniqueViolation).SetType(gin.ErrorTypePublic)
			default:
				_, file, line, _ := runtime.Caller(0)
				ctxAddGtInternalError(errMessage, file, line, err, ctx)
			}
			return
		}
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError(errMessage, file, line, err, ctx)
		return
	}

	logging.LogObjectEvent(
		ctx.FullPath(),
		ctx.ClientIP(),
		logging.ObjectEventCreate,
		nil,
		&user,
		nil,
		logging.ObjectEventSubUser,
	)
	ctx.JSON(http.StatusCreated, gin.H{"status": "created", "user": user})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	tokenUserId, _, _, err := util.GetTokenVariables(ctx)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("failed to get claims from jwt", file, line, err, ctx)
		return
	}

	reqUser, err := uc.db.GetUserById(ctx, tokenUserId)
	if err != nil {
		ctx.Error(
			gterrors.NewGtAuthError(
				gterrors.GtAuthErrorReasonTokenInvalid,
				err,
			),
		).SetType(util.GetGinErrorType())
		return
	}

	userIDToUpdate := ctx.Param("id")

	if userIDToUpdate != reqUser.ID && !reqUser.IsAdmin {
		logging.LogSecurityEvent(
			logging.SecurityScoreHigh,
			logging.SecurityEventForbiddenAction,
			ctx.FullPath(),
			userIDToUpdate,
			reqUser.Username,
		)
		ctx.Error(gterrors.ErrForbidden).SetType(gin.ErrorTypePublic)
		return
	}
	var payload *schemas.UpdateUser
	if ok := shouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	}
	isUsernameValid, err := util.ValidateUsername(payload.Username)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("failed to validate new username", file, line, err, ctx)
		return
	} else if !isUsernameValid {
		ctx.Error(gterrors.ErrUsernameUnsatisfied).SetType(gin.ErrorTypePublic)
		return
	}

	if !reqUser.IsAdmin && *payload.IsAdmin {
		logging.LogSecurityEvent(
			logging.SecurityScoreHigh,
			logging.SecurityEventForbiddenAction,
			ctx.FullPath(),
			userIDToUpdate,
			reqUser.Username,
		)
		ctx.Error(gterrors.ErrForbidden).SetType(gin.ErrorTypePublic)
		return
	}

	var oldUser *db.User
	if userIDToUpdate != reqUser.ID {
		userFromDB, err := uc.db.GetUserById(ctx, userIDToUpdate)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				ctx.Error(gterrors.ErrNotFound).SetType(gin.ErrorTypePublic)
				return
			}
			_, file, line, _ := runtime.Caller(0)
			ctxAddGtInternalError("could not get user from db", file, line, err, ctx)
			return
		}
		oldUser = &userFromDB
	} else {
		oldUser = &reqUser
	}
	if oldUser.Username == payload.Username && oldUser.IsAdmin == *payload.IsAdmin {
		logging.LogObjectEvent(
			ctx.FullPath(),
			ctx.ClientIP(),
			logging.ObjectEventUpdate,
			&reqUser,
			&oldUser,
			&oldUser,
			logging.ObjectEventSubUser,
		)
		ctx.JSON(http.StatusNoContent, gin.H{})
		return
	}

	args := &db.UpdateUserParams{
		ID: userIDToUpdate,
		Username: payload.Username,
		IsAdmin: *payload.IsAdmin,
	}

	updatedUser, err := uc.db.UpdateUser(ctx, *args)
	if err != nil {
		var pgErr *pgconn.PgError
		errMessage := "failed to update user"
		if errors.As(err, &pgErr) {
			fmt.Println("pgErr: ",pgErr)
			switch pgErr.Code {
			case "23505":
				ctx.Error(gterrors.ErrUniqueViolation).SetType(gin.ErrorTypePublic)
			default:
				_, file, line, _ := runtime.Caller(0)
				ctxAddGtInternalError(errMessage, file, line, err, ctx)
			}
			return
		}
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError(errMessage, file, line, err, ctx)
		return
	}

	if oldUser.IsAdmin != updatedUser.IsAdmin {
		if err := uc.db.DeleteJwtTokensByUserId(ctx, updatedUser.ID); err != nil {
			var pgErr *pgconn.PgError
			errMessage := "failed to delete old jwts"
			if errors.As(err, &pgErr) {
				fmt.Println("pgErr: ",pgErr)
				switch pgErr.Code {
				default:
					_, file, line, _ := runtime.Caller(0)
					logging.LogError(
						fmt.Errorf("%s: %w", errMessage, err),
						fmt.Sprintf("%v: %d", file, line),
						err.Error(),
					)
				}
				return
			}
			_, file, line, _ := runtime.Caller(0)
			logging.LogError(
				fmt.Errorf("%s: %w", errMessage, err),
				fmt.Sprintf("%v: %d", file, line),
				err.Error(),
			)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"user": updatedUser,
	})
}

// Controller for deleting users
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	tokenUserId, tokenUserName, _, err := util.GetTokenVariables(ctx)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("failed to get claims from jwt", file, line, err, ctx)
		return
	}

	reqUser, err := uc.db.GetUserById(ctx, tokenUserId)
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
		).SetType(util.GetGinErrorType())
		return
	}

	userIDToDelete := ctx.Param("id")
	if userIDToDelete != reqUser.ID && !reqUser.IsAdmin {
		logging.LogSecurityEvent(
			logging.SecurityScoreMedium,
			logging.SecurityEventForbiddenAction,
			ctx.FullPath(),
			userIDToDelete,
			reqUser.Username,
		)
		ctx.Error(gterrors.ErrForbidden).SetType(util.GetGinErrorType())
		return
	}

	rows, err := uc.db.DeleteUser(ctx, userIDToDelete)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctxAddGtInternalError("could not delete user", file, line, err, ctx)
		return
	}
	if rows == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "user-not-removed"})
		return
	}

	logging.LogObjectEvent(
		ctx.FullPath(),
		ctx.ClientIP(),
		logging.ObjectEventDelete,
		&reqUser,
		"deleted",
		userIDToDelete,
		logging.ObjectEventSubUser,
	)
	ctx.JSON(http.StatusNoContent, gin.H{})
}