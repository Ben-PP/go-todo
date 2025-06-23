package controllers

import (
	"context"
	db "go-todo/db/sqlc"
	"go-todo/schemas"
	"go-todo/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	db *db.Queries
	ctx context.Context
}

func NewAuthController(db *db.Queries, ctx context.Context) *AuthController {
	return &AuthController{db: db, ctx: ctx}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var payload *schemas.Login
	if err:= ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "malformed-body",
			"detail": err.Error(),
		})
		return
	}

	username := payload.Username
	password := payload.Password

	user, err := ac.db.GetUserByUsername(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-credentials",
			"detail": err.Error(),
		})
		return
	}

	if pwdCorrect := util.VerifyPassword(password, user.PasswordHash); !pwdCorrect {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-credentials",
		})
		return
	}

	accessToken, err := util.GenerateAccessToken(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-credentials",
			"detail": err.Error(),
		})
		return
	}
	refreshToken, err := util.GenerateRefreshToken(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-credentials",
			"detail": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}