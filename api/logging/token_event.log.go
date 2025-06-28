package logging

import (
	"go-todo/util"
	"log/slog"
)

type TokenEventType int
const (
	TokenEventTypeUse	TokenEventType	= iota
	TokenEventtypeCreate
)
func (t TokenEventType)String() string {
	switch t {
	case TokenEventtypeCreate:
		return "token:create"
	case TokenEventTypeUse:
		return "token:use"
	}
	return "unknown"
}

func LogTokenEvent(
	success bool,
	targetPath string,
	eventType TokenEventType,
	srcIp string,
	token *util.MyCustomClaims,
	) {
		if token != nil {
			LogAuditEvent(
				success,
				targetPath,
				srcIp,
				eventType.String(),
				slog.Group(
					"token",
					slog.String("sub", token.Subject),
					slog.Bool("is_admin", token.IsAdmin),
					slog.String("jti", token.ID),
					slog.String("issuer", token.Issuer),
					slog.String("issued_at", token.IssuedAt.String()),
					slog.String("family", token.Family),
					slog.String("expires_at", token.ExpiresAt.String()),
				),
			)
		} else {
			LogAuditEvent(
				success,
				targetPath,
				srcIp,
				eventType.String(),
			)
		}
}