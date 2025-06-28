package logging

import (
	db "go-todo/db/sqlc"
	"log/slog"
)

type SessionEventType int
const(
	SessionEventTypeLogin		SessionEventType = iota
	SessionEventTypeLogout
	SessionEventTypeRefresh
)

func (s SessionEventType)String() string {
	switch s {
	case SessionEventTypeLogin:
		return "session:login"
	case SessionEventTypeLogout:
		return "session:logout"
	case SessionEventTypeRefresh:
		return "session:refresh"
	}
	return "unknown"
}

func LogSessionEvent(
	success bool,
	targetPath string,
	targetUser *db.User,
	eventType SessionEventType,
	srcIp string,
) {
	LogAuditEvent(
		success,
		targetPath,
		srcIp,
		eventType.String(),
		slog.Group(
			"target",
			slog.String("username", targetUser.Username),
		),
	)
}