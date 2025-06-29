package logging

import "log/slog"

type SecurityScore int
const (
	SecurityScoreLow		SecurityScore = 1
	SecurityScoreMedium		SecurityScore = 5
	SecurityScoreHigh		SecurityScore = 10
	SecurityScoreCritical	SecurityScore = 15
)
type SecurityEventName int

const (
	SecurityEventJwtReuse		SecurityEventName = iota
	SecurityEventJwtUnknown
	SecurityEventInvalidTokenSignature
	SecurityEventLoginToUnknownUsername
	SecurityEventFailedLogin
)

func (s SecurityEventName)String() string {
	switch s {
	case SecurityEventFailedLogin:
		return "failed-login"
	case SecurityEventInvalidTokenSignature:
		return "invalid-signature-token-use"
	case SecurityEventLoginToUnknownUsername:
		return "login-to-unknown-username"
	case SecurityEventJwtReuse:
		return "jwt-reuse"
	case SecurityEventJwtUnknown:
		return "jwt-unknown"
	}
	return "unknown"	
}

func LogSecurityEvent(
	score SecurityScore,
	eventName SecurityEventName,
	target string,
	) {
	log(
		slog.LevelInfo,
		"Security event has happened",
		"security",
		slog.Int("score", int(score)),
		slog.Group(
			"event",
			slog.String("name", eventName.String()),
			slog.String("target", target),
		),
	)
}