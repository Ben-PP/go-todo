package controllers

import (
	"context"
	"errors"
	"fmt"
	db "go-todo/db/sqlc"
	"go-todo/schemas"
	"go-todo/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "malformed-body",
			"detail": err.Error(),
		})
		return
	}

	users, err := uc.db.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
			"detail": err.Error(),
		})
		return
	}

	makeAdmin := false
	if len(users) == 0 {
		makeAdmin = true
	}

	isPasswdValid, err := util.ValidatePassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "internal-server-error"})
		return
	} else if !isPasswdValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "password-criteria-unmet"})
		return
	}

	isUsernameValid, err := util.ValidateUsername(payload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "internal-server-error"})
		return
	} else if !isUsernameValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "username-criteria-unmet"})
		return
	}

	userUUID := uuid.New()
	passwd := payload.Password
	passwdHash,err := util.HashPassword(passwd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "unable-to-hash-password",
			"detail": err.Error(),
		})
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
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusConflict, gin.H{
					"status": "unique-violation",
					"detail": pgErr.Detail,
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": "unable-to-create-user",
					"detail": pgErr.Error(),
				})
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "unable-to-create-user",
			"detail": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "created", "user": user})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	tokenUserId, _, _, err := util.GetTokenVariables(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "invalid-token"})
		return
	}

	reqUser, err := uc.db.GetUserById(ctx, tokenUserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "invalid-token"})
		return
	}

	userIDToUpdate := ctx.Param("id")

	if userIDToUpdate != reqUser.ID && !reqUser.IsAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "forbidden",
			"detail": "Only admin can modify other users.",
		})
		return
	}
	var payload *schemas.UpdateUser
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "malformed-body",
			"detail": err.Error(),
		})
		return
	}
	isUsernameValid, err := util.ValidateUsername(payload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "internal-server-error"})
		return
	} else if !isUsernameValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "username-criteria-unmet"})
		return
	}

	if !reqUser.IsAdmin && *payload.IsAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "forbidden",
			"detail": "Only admin can promote user to admin.",
		})
		return
	}

	var oldUser *db.User
	if userIDToUpdate != reqUser.ID {
		userFromDB, err := uc.db.GetUserById(ctx, userIDToUpdate)
		if err != nil {
			if err.Error() == "no rows in result set" {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status": "user-not-found",
					"detail": "User "+userIDToUpdate+" was not found",
				})
				return
			}
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				default:
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"status": "internal-server-error",
						"detail": pgErr.Error(),
					})
				}
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "internal-server-error", "detail": err.Error()})
			return
			// TODO Handle user not found. This is not server error and should return a proper response.
		}
		oldUser = &userFromDB
	} else {
		oldUser = &reqUser
	}
	if oldUser.Username == payload.Username && oldUser.IsAdmin == *payload.IsAdmin {
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
		if errors.As(err, &pgErr) {
			fmt.Println("pgErr: ",pgErr)
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusConflict, gin.H{
					"status": "unique-violation",
					"detail": pgErr.Detail,
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status": "internal-server-error",
					"detail": pgErr.Error(),
				})
			}
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
			"detail": err.Error(),
		})
		return
	}

	if oldUser.IsAdmin != updatedUser.IsAdmin {
		if err := uc.db.DeleteJwtTokensByUserId(ctx, updatedUser.ID); err != nil {
			// TODO Log this as this would be bad if ever happened or user has not logged in once yet
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				fmt.Println("pgErr: ",pgErr)
				switch pgErr.Code {
				default:
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"status": "internal-server-error",
						"detail": pgErr.Error(),
					})
				}
				return
			}
			fmt.Println("Error: ",err.Error())
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"user": updatedUser,
	})
}

// Controller for deleting users
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	tokenUserId, _, _, err := util.GetTokenVariables(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "invalid-token"})
		return
	}

	reqUser, err := uc.db.GetUserById(ctx, tokenUserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "invalid-token"})
		return
	}

	userIDToDelete := ctx.Param("id")
	if userIDToDelete != reqUser.ID && !reqUser.IsAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "forbidden",
			"detail": "Only admins can delete other users.",
		})
		return
	}

	rows, err := uc.db.DeleteUser(ctx, userIDToDelete)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "internal-server-error"})
		return
	}
	if rows == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "user-not-removed"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}