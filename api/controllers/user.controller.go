package controllers

import (
	"context"
	"errors"
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

	userUUID := uuid.New()
	// TODO Check password requirements
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
		IsAdmin: true,
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