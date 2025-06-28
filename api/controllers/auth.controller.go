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
	"github.com/jackc/pgx/v5"
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

func logTokenEventUse(success bool, token *util.MyCustomClaims, c *gin.Context) {
	logging.LogTokenEvent(success, c.FullPath(), logging.TokenEventTypeUse,	c.RemoteIP(), token)
}

func logTokenCreations(claims []*util.MyCustomClaims, c *gin.Context) {
	for _, claims := range claims {
		logging.LogTokenEvent(
			true,
			c.FullPath(),
			logging.TokenEventtypeCreate,
			c.RemoteIP(),
			claims,
		)
	}
}

func generateTokens(family string, user db.User) (
	refreshToken string,
	refreshClaims *util.MyCustomClaims,
	accessToken string,
	accessClaims *util.MyCustomClaims,
	err error,
	) {
		refreshToken, refreshClaims, err = util.GenerateRefreshToken(
		user.Username,
		user.ID,
		user.IsAdmin,
		family,
	)
	if err != nil {
		return "", nil, "", nil, err
	}
	accessToken, accessClaims, err = util.GenerateAccessToken(
		user.Username,
		user.ID,
		user.IsAdmin,
	)
	if err != nil {

		return "", nil, "", nil, err
	}
	return
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
		logTokenEventUse(false, decodedRefreshToken, ctx)
		ctx.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	dbToken, err := ac.db.GetJwtTokenByJti(ctx, decodedRefreshToken.ID)
	if err != nil {
		logTokenEventUse(false, decodedRefreshToken, ctx)
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
		logTokenEventUse(false, decodedRefreshToken, ctx)
		logging.LogSecurityEvent(
			logging.SecurityScoreCritical,
			logging.SecurityEventRefreshTokenReuse,
			fmt.Sprintf("value=%v,type=jwt",decodedRefreshToken.ID),
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
		logTokenEventUse(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)

		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrDatabaseMeta(
			fmt.Sprintf("%v: %d", file, line),
			fmt.Sprintf("User '%v' was not found.", dbToken.UserID),
			401,
		))
		return
	}

	logSessionRefresh := func (success bool){logging.LogSessionEvent(
		success,
		ctx.FullPath(),
		&user,
		logging.SessionEventTypeRefresh,
		ctx.RemoteIP(),
	)}
	
	refreshToken, refreshClaims, accessToken, accessClaims, err := generateTokens(
		decodedRefreshToken.Family,
		user,
	)
	if err != nil {
		logSessionRefresh(false)
		logTokenEventUse(false, decodedRefreshToken, ctx)
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
		logTokenEventUse(false, decodedRefreshToken, ctx)
		logSessionRefresh(false)
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrInternalMeta(
			fmt.Sprintf("%v: %d", file, line),
			err.Error(),
		))
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
		logTokenEventUse(false, decodedRefreshToken, ctx)
		logSessionRefresh(false)
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

	logSessionRefresh(true)
	logTokenCreations([]*util.MyCustomClaims{refreshClaims,accessClaims}, ctx)
	logTokenEventUse(true, decodedRefreshToken, ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var payload *schemas.Login
	if err:= ctx.ShouldBindBodyWithJSON(&payload); err != nil {
		ctx.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	username := payload.Username
	password := payload.Password

	user, err := ac.db.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logging.LogSecurityEvent(
				logging.SecurityScoreLow,
				logging.SecurityEventLoginToInvalidUsername,
				fmt.Sprintf("value=%v,type=login",username),
			)
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "invalid-credentials"})
			return
		}
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrInternalMeta(
			fmt.Sprintf("%v: %d", file, line),
			err.Error(),
		))
		return
	}

	if pwdCorrect := util.VerifyPassword(password, user.PasswordHash); !pwdCorrect {
		logging.LogSecurityEvent(
			logging.SecurityScoreLow,
			logging.SecurityEventFailedLogin,
			fmt.Sprintf("value=%v,type=login",username),
		)
		logging.LogSessionEvent(
			false,
			ctx.FullPath(),
			&user,
			logging.SessionEventTypeLogin,
			ctx.RemoteIP(),
		)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "invalid-credentials",
		})
		return
	}

	refreshToken, refreshClaims, accessToken, accessClaims, err := generateTokens(
		"",
		user,
	)
	if err != nil {
		logTokenEventUse(false, &util.MyCustomClaims{}, ctx)
		_, file, line, _ := runtime.Caller(0)

		ctx.Error(err).SetType(gin.ErrorTypePrivate).
		SetMeta(*util.NewErrInternalMeta(
			fmt.Sprintf("%v: %d", file, line),
			err.Error(),
		))
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
		logTokenEventUse(false, &util.MyCustomClaims{}, ctx)
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

	logging.LogSessionEvent(
		true,
		ctx.FullPath(),
		&user,
		logging.SessionEventTypeLogin,
		ctx.RemoteIP(),
	)
	logTokenCreations([]*util.MyCustomClaims{refreshClaims,accessClaims}, ctx)
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