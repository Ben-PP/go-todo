package controllers

import (
	"context"
	"errors"
	"fmt"
	db "go-todo/db/sqlc"
	"go-todo/logging"
	"go-todo/schemas"
	"go-todo/util"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthController struct {
	db *db.Queries
	ctx context.Context
}

func NewAuthController(db *db.Queries, ctx context.Context) *AuthController {
	return &AuthController{db: db, ctx: ctx}
}

func logTokenUsage(success bool, token *util.MyCustomClaims, c *gin.Context) {
	logging.LogTokenUsage(success, c.FullPath(), "use",	"refresh", c.RemoteIP(), token,)
}

func (ac *AuthController) Refresh(ctx *gin.Context) {
	var payload *schemas.Refresh
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	refreshToken := payload.RefreshToken

	decodedRefreshToken, err := util.DecodeRefreshToken(refreshToken)
	if err != nil {
		logTokenUsage(false, decodedRefreshToken, ctx)
		ctx.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	dbToken, err := ac.db.GetJwtTokenByJti(ctx, decodedRefreshToken.ID)
	if err != nil {
		logTokenUsage(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)
		
		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrDatabaseMeta(
			fmt.Sprintf("%v: %d", file, line),
			fmt.Sprintf("Queried for jwt token: %v", decodedRefreshToken.ID),
			401,
		))
		return
	}
	
	if dbToken.IsUsed {
		logTokenUsage(false, decodedRefreshToken, ctx)
		logging.LogSecurityEvent(
			logging.SecurityScoreHigh,
			logging.SecurityEventRefreshTokenReuse,
		)
		if err := ac.db.DeleteJwtTokenByFamily(ctx, dbToken.Family); err != nil {
			_, file, line, _ := runtime.Caller(0)
			logging.LogError(
				err,
				fmt.Sprintf("%v: %d", file, line),
				fmt.Sprintf("Failed to delete jwt family '%v'", dbToken.Family),
			)
		}
		ctx.JSON(401, gin.H{"status": "unauthorized"})
		return
	}
	
	user, err := ac.db.GetUserById(ctx, dbToken.UserID)
	if err != nil {
		logTokenUsage(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)

		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrDatabaseMeta(
			fmt.Sprintf("%v: %d", file, line),
			fmt.Sprintf("User '%v' was not found.", dbToken.UserID),
			401,
		))
		return
	}
	
	newRefreshToken, refreshClaims, err := util.GenerateRefreshToken(
		user.Username,
		user.ID,
		user.IsAdmin,
		decodedRefreshToken.Family,
	)
	if err != nil {
		logTokenUsage(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)

		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrInternalMeta(
			fmt.Sprintf("%v: %d", file, line),
			err.Error(),
		))
		return
	}
	newAccessToken, _, err := util.GenerateAccessToken(
		user.Username,
		user.ID,
		user.IsAdmin,
	)
	if err != nil {
		logTokenUsage(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrInternalMeta(
			fmt.Sprintf("%v: %d", file, line),
			err.Error(),
		))
		return
	}

	// Mark the token as used.
	if err := ac.db.UseJwtToken(ctx, dbToken.Jti); err != nil {
		logTokenUsage(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrInternalMeta(
			fmt.Sprintf("%v: %d", file, line),
			err.Error(),
		))
	}

	args := &db.CreateJwtTokenParams{
		Jti: refreshClaims.ID,
		UserID: refreshClaims.Subject,
		Family: refreshClaims.Family,
		CreatedAt: pgtype.Timestamp{Time: refreshClaims.IssuedAt.Time,Valid: true},
		ExpiresAt: pgtype.Timestamp{Time: refreshClaims.ExpiresAt.Time, Valid: true},
	}
	
	if err := ac.db.CreateJwtToken(ctx, *args); err != nil {
		logTokenUsage(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			ctx.Error(pgErr).SetType(gin.ErrorTypePrivate).
			SetMeta(*util.NewErrDatabaseMeta(
				fmt.Sprintf("%v: %d", file, line),
				fmt.Sprintf("Unable to insert token '%v'", args.Jti),
				500,
			))
		} else {		
			ctx.Error(err).SetType(gin.ErrorTypePrivate).
			SetMeta(*util.NewErrDatabaseMeta(
				fmt.Sprintf("%v: %d", file, line),
				fmt.Sprintf("Unable to insert token '%v'", args.Jti),
				500,
			))
		}
		return
	}

	logTokenUsage(true, decodedRefreshToken, ctx)
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

func (ac *AuthController) Logout(ctx *gin.Context) {
	var payload *schemas.Refresh
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "malformed-body",
			"detail": err.Error(),
		})
		return
	}

	refreshToken := payload.RefreshToken

	claims, err := util.DecodeRefreshToken(refreshToken)
	if err != nil {
		// TODO Log
		// Providing access token instead of invalid token leads to badnes
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-token",
		})
		return
	}

	if err := ac.db.DeleteJwtTokenByFamily(ctx,claims.Family); err != nil {
		// TODO Log
		// Running here might cause incomplete logout
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (ac *AuthController) UpdatePassword(ctx *gin.Context) {
	userID,_,_, err := util.GetTokenVariables(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "invalid-token"})
		return
	}

	var payload *schemas.UpdatePassword
	if err := ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "malformed-body",
			"detail": err.Error(),
		})
		return
	}
	
	isPasswdValid, err := util.ValidatePassword(payload.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "internal-server-error"})
		return
	} else if !isPasswdValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "password-criteria-unmet"})
		return
	}

	user, err := ac.db.GetUserById(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "unauthorized",
			"detail": "Tokens user does not exists.",
		})
		return
	}

	if !util.VerifyPassword(payload.OldPassword, user.PasswordHash) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-password",
			"detail": "Old password does not match",
		})
		return
	}

	newPasswordHash, err := util.HashPassword(payload.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
		})
		return
	}

	args := &db.UpdateUserPasswordParams{
		PasswordHash: newPasswordHash,
		ID: user.ID,
	}

	if err := ac.db.UpdateUserPassword(ctx, *args); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
		})
	}

	accessToken, _, err := util.GenerateAccessToken(
		user.Username,
		user.ID,
		user.IsAdmin,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
			"detail": err.Error(),
		})
		return
	}

	refreshArgs := &db.CreateJwtTokenParams{
		Jti: refreshClaims.ID,
		UserID: refreshClaims.Subject,
		Family: refreshClaims.Family,
		CreatedAt: pgtype.Timestamp{Time: refreshClaims.IssuedAt.Time,Valid: true},
		ExpiresAt: pgtype.Timestamp{Time: refreshClaims.ExpiresAt.Time, Valid: true},
	}
	
	if err := ac.db.CreateJwtToken(ctx, *refreshArgs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "internal-server-error",
			"detail": err.Error(),
		})
		return
	}

	deleteArgs := &db.DeleteJwtTokenByUserIdExcludeFamilyParams{
		UserID: userID,
		Family: refreshClaims.Family,
	}

	if err := ac.db.DeleteJwtTokenByUserIdExcludeFamily(ctx, *deleteArgs); err != nil {
		// TODO Log
		// Running here might cause incomplete logout
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

// TODO Add ResetPassword (Requires email to be implemented)