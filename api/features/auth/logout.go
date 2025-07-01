package auth

import (
	"errors"
	"fmt"
	"go-todo/gterrors"
	"go-todo/logging"
	"go-todo/schemas"
	"go-todo/util/jwt"
	"go-todo/util/mycontext"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

func (ac *AuthController) Logout(ctx *gin.Context) {
	// TODO Add access jwt to redis blacklist
	// TODO Validate that access tokens user matches with the refresh tokens user
	var payload *schemas.Refresh
	if ok := mycontext.ShouldBindBodyWithJSON(&payload, ctx); !ok {
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
				).SetType(gterrors.GetGinErrorType())
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
		mycontext.CtxAddGtInternalError("", file, line, errIfNil, ctx)
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