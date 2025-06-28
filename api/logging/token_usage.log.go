package logging

import (
	"go-todo/util"
	"log/slog"
)

func LogTokenUsage(
	success bool,
	targetPath string,
	eventType string,
	usedAsType string,
	srcIp string,
	token *util.MyCustomClaims,
	) {
		origin := slog.Group(
			"origin",
			slog.String("ip", srcIp),
		)
		if token != nil {
			log(
				slog.LevelInfo,
				"Access event to the application",
				"tokens",
				origin,
				slog.Group(
					"action",
					slog.String("type", eventType),
					slog.Bool("success", success),
					slog.String("target_path", targetPath),
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
				origin,
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