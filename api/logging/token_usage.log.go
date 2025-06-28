package logging

import (
	"go-todo/util"
	"log/slog"
)

func LogTokenUsage(
	success bool,
	eventType string,
	usedAsType string,
	token *util.MyCustomClaims,
	/*tokenID string,
	tokenExpiry string,*/
	) {
		if token != nil {
			log(
				slog.LevelInfo,
				"Access event to the application",
				"tokens",
				slog.Group(
					"action",
					slog.String("type", eventType),
					slog.Bool("success", success),
					slog.String("used_as_type", usedAsType),
				),
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
			log(
				slog.LevelInfo,
				"Access event to the application",
				"tokens",
				slog.Group(
					"action",
					slog.String("type", eventType),
					slog.Bool("success", success),
					slog.String("used_as_type", usedAsType),
				),
				slog.String("token", "nil"),
			)
		}
}