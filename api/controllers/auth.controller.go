package controllers

import (
	"context"
	"fmt"
	db "go-todo/db/sqlc"
	"go-todo/schemas"
	"go-todo/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthController struct {
	db *db.Queries
	ctx context.Context
}

func NewAuthController(db *db.Queries, ctx context.Context) *AuthController {
	return &AuthController{db: db, ctx: ctx}
}

func invalidTokenResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"status": "invalid-token",
	})
}

func (ac *AuthController) Refresh(ctx *gin.Context) {
	var payload *schemas.Refresh
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "malformed-body",
			"detail": err.Error(),
		})
	}

	refreshToken := payload.RefreshToken

	decodedRefreshToken, err := util.DecodeRefreshToken(refreshToken)
	if err != nil {
		invalidTokenResponse(ctx)
		return
	}

	dbToken, err := ac.db.GetJwtTokenByJti(ctx, decodedRefreshToken.ID)
	if err != nil {
		invalidTokenResponse(ctx)
		fmt.Println(err.Error())
		return
	}

	if dbToken.IsUsed {
		if err := ac.db.DeleteJwtTokenByFamily(ctx, dbToken.Family); err != nil {
			// TODO Log error
		}
		invalidTokenResponse(ctx)
		return
	}

	newRefreshToken, refreshClaims, err := util.GenerateRefreshToken(
		decodedRefreshToken.UserName,
		decodedRefreshToken.Subject,
		decodedRefreshToken.IsAdmin,
		decodedRefreshToken.Family,
	)
	if err != nil {
		invalidTokenResponse(ctx)
		return
	}
	newAccessToken, _, err := util.GenerateAccessToken(
		decodedRefreshToken.UserName,
		decodedRefreshToken.Subject,
		decodedRefreshToken.IsAdmin,
	)
	if err != nil {
		invalidTokenResponse(ctx)
		return
	}

	ac.db.UseJwtToken(ctx, dbToken.Jti)

	args := &db.CreateJwtTokenParams{
		Jti: refreshClaims.ID,
		UserID: refreshClaims.Subject,
		Family: refreshClaims.Family,
		CreatedAt: pgtype.Timestamp{Time: refreshClaims.IssuedAt.Time,Valid: true},
		ExpiresAt: pgtype.Timestamp{Time: refreshClaims.ExpiresAt.Time, Valid: true},
	}
	
	if err := ac.db.CreateJwtToken(ctx, *args); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
			"detail": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"access_token": newAccessToken,
		"refresh_token": newRefreshToken,
	})
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

	accessToken, _, err := util.GenerateAccessToken(
		user.Username,
		user.ID,
		user.IsAdmin,
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-credentials",
			"detail": err.Error(),
		})
		return
	}
	refreshToken, refreshClaims, err := util.GenerateRefreshToken(
		user.Username,
		user.ID,
		user.IsAdmin,
		"",
	)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-credentials",
			"detail": err.Error(),
		})
		return
	}

	args := &db.CreateJwtTokenParams{
		Jti: refreshClaims.ID,
		UserID: refreshClaims.Subject,
		Family: refreshClaims.Family,
		CreatedAt: pgtype.Timestamp{Time: refreshClaims.IssuedAt.Time,Valid: true},
		ExpiresAt: pgtype.Timestamp{Time: refreshClaims.ExpiresAt.Time, Valid: true},
	}
	
	if err := ac.db.CreateJwtToken(ctx, *args); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
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