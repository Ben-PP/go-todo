package logging

import "log/slog"

type SecurityScore int
const (
	SecurityScoreLow		SecurityScore = 1
	SecurityScoreMedium	SecurityScore = 2
	SecurityScoreHigh	SecurityScore = 3
)
type SecurityEventName int

const (
	SecurityEventRefreshTokenReuse		SecurityEventName = iota
	SecurityEventInvalidTokenSignature
)

func (s SecurityEventName)String() string {
	switch s {
	case SecurityEventRefreshTokenReuse:
		return "refresh-token-reuse"
	case SecurityEventInvalidTokenSignature:
		return "invalid-signature-token-use"
	}
	return "unknown"	
}

func LogSecurityEvent(score SecurityScore, eventName SecurityEventName) {
	log(
		slog.LevelInfo,
		"Security event has happened",
		"security",
		slog.Int("score", int(score)),
		slog.Group(
			"event",
			slog.String("name", eventName.String()),
		),
	)
}