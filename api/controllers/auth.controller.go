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
	"go-todo/util/jwt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthController struct {
	db *db.Queries
	ctx context.Context
}

func NewAuthController(db *db.Queries, ctx context.Context) *AuthController {
	return &AuthController{db: db, ctx: ctx}
}

func logTokenEventUse(success bool, token *jwt.GtClaims, c *gin.Context) {
	logging.LogTokenEvent(success, c.FullPath(), logging.TokenEventTypeUse,	c.RemoteIP(), token)
}

func logTokenCreations(claims []*jwt.GtClaims, c *gin.Context) {
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
	refreshClaims *jwt.GtClaims,
	accessToken string,
	accessClaims *jwt.GtClaims,
	err error,
	) {
		refreshToken, refreshClaims, err = jwt.GenerateRefreshJwt(
		user.Username,
		user.ID,
		user.IsAdmin,
		family,
	)
	if err != nil {
		return "", nil, "", nil, err
	}
	accessToken, accessClaims, err = jwt.GenerateAccessJwt(
		user.Username,
		user.ID,
		user.IsAdmin,
	)
	if err != nil {

		return "", nil, "", nil, err
	}
	return
}

func createGtInternalError(msg, file string, line int, err error, c *gin.Context) {
	c.Error(
		gterrors.NewGtInternalError(
			fmt.Errorf("%v: %w", msg, err),
			fmt.Sprintf("%v: %d", file, line),
			500,
		),
	).SetType(util.GetGinErrorType())
}

func failedToGenerateJwtError(err error, file string, line int, c *gin.Context) {
	createGtInternalError("failed to generate jwt", file, line, err, c)
}

func failedToSaveJwtToDbError(err error, file string, line int, c *gin.Context) {
	createGtInternalError("failed to save jwt to db", file, line, err, c)
}
			
func (ac *AuthController) Refresh(ctx *gin.Context) {
	var payload *schemas.Refresh
	if ok := shouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	}

	refreshToken := payload.RefreshToken

	decodedRefreshToken, err := jwt.DecodeRefreshToken(refreshToken)
	if err != nil {
		logTokenEventUse(false, decodedRefreshToken, ctx)
		var jwtErr *jwt.JwtDecodeError
		if errors.As(err, &jwtErr) {
			reason := gterrors.GtAuthErrorReasonInternalError
			switch jwtErr.Reason {
			case jwt.JwtErrorReasonExpired:
				reason = gterrors.GtAuthErrorReasonExpired
			case jwt.JwtErrorReasonInvalidSignature:
				reason = gterrors.GtAuthErrorReasonInvalidSignature
			case jwt.JwtErrorReasonTokenMalformed:
				reason = gterrors.GtAuthErrorReasonTokenInvalid
			case jwt.JwtErrorReasonUnhandled:
				reason = gterrors.GtAuthErrorReasonInternalError
			}

			ctx.Error(gterrors.NewGtAuthError(reason, err)).SetType(util.GetGinErrorType())
			return
		}
		// Should never get to here
		ctx.Error(gterrors.ErrShouldNotHappen)
		return
	}

	dbToken, err := ac.db.GetJwtTokenByJti(ctx, decodedRefreshToken.ID)
	if err != nil {
		logTokenEventUse(false, decodedRefreshToken, ctx)
		logging.LogSecurityEvent(
			logging.SecurityScoreMedium,
			logging.SecurityEventJwtUnknown,
			decodedRefreshToken.ID,
		)
		ginType := util.GetGinErrorType()
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.Error(
				gterrors.NewGtAuthError(gterrors.GtAuthErrorReasonTokenInvalid, err),
			).SetType(ginType)
			return
		}

		ctx.Error(gterrors.NewGtAuthError(
			gterrors.GtAuthErrorReasonInternalError,
			fmt.Errorf("failed to get token from db: %w", err),
		)).SetType(ginType)
		return
	}
	
	if dbToken.IsUsed {
		logTokenEventUse(false, decodedRefreshToken, ctx)
		logging.LogSecurityEvent(
			logging.SecurityScoreCritical,
			logging.SecurityEventJwtReuse,
			fmt.Sprintf("value=%v,type=jwt",decodedRefreshToken.ID),
		)
		ginType := util.GetGinErrorType()
		if rows, err := ac.db.DeleteJwtTokenByFamily(ctx, dbToken.Family);
		err != nil || rows == 0 {
			_, file, line, _ := runtime.Caller(0)
			errIfNil := fmt.Errorf("failed to delete jwt family: %w", err)
			if err == nil {
				errIfNil = fmt.Errorf("failed to delete jwt family: %v", dbToken.Family)
			}
			fileFull := fmt.Sprintf("%v: %d", file, line)
			ctx.Error(gterrors.NewGtInternalError(errIfNil,	fileFull, 500)).SetType(ginType)
			return
		}
		ctx.Error(
			gterrors.NewGtAuthError(
				gterrors.GtAuthErrorReasonTokenReuse,
				gterrors.ErrJwtRefreshReuse,
			),
		).SetType(ginType)
		return
	}
	
	// This should always succeed if db works correctly as dbToken has to have
	// userID. Errors are system failures.
	user, err := ac.db.GetUserById(ctx, dbToken.UserID)
	if err != nil {
		logTokenEventUse(false, decodedRefreshToken, ctx)
		_, file, line, _ := runtime.Caller(0)

		ctx.Error(gterrors.NewGtInternalError(
			fmt.Errorf("failed to get user from db: %w", err),
			fmt.Sprintf("%v: %d", file, line),
			500,	
		)).SetType(util.GetGinErrorType())
		return
	}

	logSessionRefresh := func (success bool){logging.LogSessionEvent(
		success,
		ctx.FullPath(),
		user.Username,
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
		failedToGenerateJwtError(err, file, line, ctx)
		return
	}

	// Mark the token as used.
	if err := ac.db.UseJwtToken(ctx, dbToken.Jti); err != nil {
		logTokenEventUse(false, decodedRefreshToken, ctx)
		logSessionRefresh(false)
		_, file, line, _ := runtime.Caller(0)

		ctx.Error(
			gterrors.NewGtInternalError(
				err,
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		).SetType(util.GetGinErrorType())
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
		failedToSaveJwtToDbError(err, file, line, ctx)
		return
	}

	logSessionRefresh(true)
	logTokenCreations([]*jwt.GtClaims{refreshClaims,accessClaims}, ctx)
	logTokenEventUse(true, decodedRefreshToken, ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var payload *schemas.Login
	if ok := shouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	}

	username := payload.Username
	password := payload.Password
	ok, err := util.ValidateUsername(username)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(
			gterrors.NewGtInternalError(
				err,
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		).SetType(util.GetGinErrorType())
		return
	} else if !ok {
		ctx.Error(
			gterrors.NewGtAuthError(
				gterrors.GtAuthErrorReasonUsernameInvalid,
				errors.New("username validation failed"),
			),
		).SetType(util.GetGinErrorType())
		return
	}

	user, err := ac.db.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logging.LogSecurityEvent(
				logging.SecurityScoreLow,
				logging.SecurityEventLoginToUnknownUsername,
				username,
			)

			ctx.Error(
				gterrors.NewGtAuthError(
					gterrors.GtAuthErrorReasonInvalidCredentials,
					err,
				),
			).SetType(util.GetGinErrorType())
			return
		}

		_, file, line, _ := runtime.Caller(0)
		ctx.Error(
			gterrors.NewGtInternalError(
				fmt.Errorf("failed to get username from db: %w", err),
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		).SetType(util.GetGinErrorType())
		return
	}

	if pwdCorrect := util.ComparePassword(password, user.PasswordHash); !pwdCorrect {
		logging.LogSecurityEvent(
			logging.SecurityScoreLow,
			logging.SecurityEventFailedLogin,
			username,
		)
		logging.LogSessionEvent(
			false,
			ctx.FullPath(),
			user.Username,
			logging.SessionEventTypeLogin,
			ctx.RemoteIP(),
		)

		ctx.Error(
			gterrors.NewGtAuthError(
				gterrors.GtAuthErrorReasonInvalidCredentials,
				errors.New("password verification failed"),
			),
		).SetType(util.GetGinErrorType())
		return
	}

	refreshToken, refreshClaims, accessToken, accessClaims, err := generateTokens(
		"",
		user,
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		failedToGenerateJwtError(err, file, line, ctx)
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
		_, file, line, _ := runtime.Caller(0)	
		failedToSaveJwtToDbError(err, file, line, ctx)
		return
	}

	logging.LogSessionEvent(
		true,
		ctx.FullPath(),
		user.Username,
		logging.SessionEventTypeLogin,
		ctx.ClientIP(),
	)
	logTokenCreations([]*jwt.GtClaims{refreshClaims,accessClaims}, ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	// TODO Add access jwt to redis blacklist
	// TODO Validate that access tokens user matches with the refresh tokens user
	var payload *schemas.Refresh
	if ok := shouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	}

	refreshToken := payload.RefreshToken

	claims, err := jwt.DecodeRefreshToken(refreshToken)
	if err != nil {
		logTokenEventUse(false, claims, ctx)
		var jwtErr *jwt.JwtDecodeError
		if errors.As(err, &jwtErr) {
			reason := gterrors.GtAuthErrorReasonInternalError
			switch jwtErr.Reason {
			case jwt.JwtErrorReasonExpired:
				reason = gterrors.GtAuthErrorReasonExpired
			case jwt.JwtErrorReasonInvalidSignature:
				reason = gterrors.GtAuthErrorReasonInvalidSignature
			case jwt.JwtErrorReasonTokenMalformed:
				reason = gterrors.GtAuthErrorReasonTokenInvalid
			case jwt.JwtErrorReasonUnhandled:
				reason = gterrors.GtAuthErrorReasonInternalError
			}

			ctx.Error(
				gterrors.NewGtAuthError(
					reason,
					errors.Join(gterrors.ErrGtLogoutFailure, err),
					),
				).SetType(util.GetGinErrorType())
			return
		}
		// Should never get to here
		ctx.Error(gterrors.ErrShouldNotHappen)
		return
	}

	if rows, err := ac.db.DeleteJwtTokenByFamily(ctx, claims.Family);
	err != nil || rows == 0 {
		_, file, line, _ := runtime.Caller(0)
		errIfNil := fmt.Errorf("failed to delete jwt family: %w", err)
		if err == nil {
			errIfNil = fmt.Errorf("failed to delete jwt family: %v", claims.Family)
		}
		fileFull := fmt.Sprintf("%v: %d", file, line)
		ctx.Error(gterrors.NewGtInternalError(errIfNil, fileFull, 500)).SetType(util.GetGinErrorType())
		return
	}
	
	logging.LogSessionEvent(
		true,
		ctx.FullPath(),
		claims.Username,
		logging.SessionEventTypeLogout,
		ctx.ClientIP(),
	)
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (ac *AuthController) UpdatePassword(ctx *gin.Context) {
	userID,_,_, err := util.GetTokenVariables(ctx)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(
			gterrors.NewGtInternalError(
				fmt.Errorf("failed to get claims from jwt: %w", err),
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		).SetType(util.GetGinErrorType())
		return
	}

	var payload *schemas.UpdatePassword
	if ok := shouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	}
	
	isPasswdValid, err := util.ValidatePassword(payload.NewPassword)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(
			gterrors.NewGtInternalError(
				fmt.Errorf("error validatig password: %w", err),
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		).SetType(util.GetGinErrorType())
		return
	} else if !isPasswdValid {
		ctx.Error(gterrors.ErrPasswordUnsatisfied).SetType(gin.ErrorTypePublic)
		return
	}

	// Should only fail if something is wrong in the server
	user, err := ac.db.GetUserById(ctx, userID)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(gterrors.NewGtInternalError(
			fmt.Errorf("failed to get user from db: %w", err),
			fmt.Sprintf("%v: %d", file, line),
			500,	
		)).SetType(util.GetGinErrorType())
		return
	}

	if !util.ComparePassword(payload.OldPassword, user.PasswordHash) {
		ctx.Error(
			gterrors.NewGtAuthError(
				gterrors.GtAuthErrorReasonInvalidCredentials,
				errors.New("provided credentials are incorrect"),
			),
		).SetType(gin.ErrorTypePublic)
		return
	}

	// TODO Check that the new password is not the old password. Maybe a loop?
	newPasswordHash, err := util.HashPassword(payload.NewPassword)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(
			gterrors.NewGtInternalError(
				fmt.Errorf("failed to hash new password: %w", err),
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		).SetType(util.GetGinErrorType())
		return
	}

	args := &db.UpdateUserPasswordParams{
		PasswordHash: newPasswordHash,
		ID: user.ID,
	}

	if err := ac.db.UpdateUserPassword(ctx, *args); err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(
			gterrors.NewGtInternalError(
				fmt.Errorf("failed to update password to db: %w", err),
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		)
		return
	}

	refreshToken, refreshClaims, accessToken, _, err := generateTokens(
		"",
		user,
	)
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		failedToGenerateJwtError(err, file, line, ctx)
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
		_, file, line, _ := runtime.Caller(0)	
		failedToSaveJwtToDbError(err, file, line, ctx)
		return
	}

	deleteArgs := &db.DeleteJwtTokenByUserIdExcludeFamilyParams{
		UserID: userID,
		Family: refreshClaims.Family,
	}

	if err := ac.db.DeleteJwtTokenByUserIdExcludeFamily(ctx, *deleteArgs); err != nil {
		_, file, line, _ := runtime.Caller(0)
		ctx.Error(
			gterrors.NewGtInternalError(
				fmt.Errorf("failed to remove old refresh jwts: %w", err),
				fmt.Sprintf("%v: %d", file, line),
				500,
			),
		).SetType(util.GetGinErrorType())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
}

// TODO Add ResetPassword (Requires email to be implemented)