package logging

import (
	db "go-todo/db/sqlc"
	"log/slog"
)

type ObjectEventSub int
const (
	ObjectEventSubUser	ObjectEventSub = iota
)

func (e ObjectEventSub) String() string {
	switch e {
	case ObjectEventSubUser:
		return "user"
	}
	return "unknown"
}

type ObjectEvent int
const (
	ObjectEventCreate	ObjectEvent = iota
	ObjectEventRead
	ObjectEventUpdate
	ObjectEventDelete
)

func (e ObjectEvent) String() string {
	switch e {
	case ObjectEventCreate:
		return "objectevent:create"
	case ObjectEventRead:
		return "objectevent:read"
	case ObjectEventUpdate:
		return "objectevent:update"
	case ObjectEventDelete:
		return "objectevent:delete"
	}
	return "objectevent:unknown"
}

// Logs object crud events. Subjects are given like &db.<type>
func LogObjectEvent(
	targetPath string,
	srcIp string,
	eventType ObjectEvent,
	actor *db.User,
	subjectCurrent any,
	subjectOld any,
	subjectType ObjectEventSub,
) {
	getActorData := func(id, username string, isAdmin bool) slog.Attr {
		return slog.Group(
			"actor",
			slog.String("id", id),
			slog.String("username", username),
			slog.Bool("is_admin", isAdmin),
		)
	}

	getSubject := func(subCur, subOld any) slog.Attr{
		var groupCurrent *slog.Attr
		var groupOld *slog.Attr

		curKey := "current"
		oldKey := "old"
		switch sc := subCur.(type) {
		case string:
			if eventType == ObjectEventDelete {
				gCur := slog.String("id", sc)
				groupCurrent = &gCur
				if so ,ok := subOld.(string); ok {
					gOld := slog.String("id", so)
					groupOld = &gOld
				}
			}
		case *db.CreateUserRow:
			gCur :=slog.Group(
				curKey,
				slog.String("id", sc.ID),
				slog.String("username", sc.Username),
				slog.Bool("is_admin", sc.IsAdmin),
			)
			groupCurrent = &gCur
			if subOld != nil {
				so := subOld.(*db.CreateUserRow)
				gOld := slog.Group(
					oldKey,
					slog.String("id", so.ID),
					slog.String("username", so.Username),
					slog.Bool("is_admin", so.IsAdmin),
				)
				groupOld = &gOld
			}
		default:
			group := slog.String(curKey, "nil")
			groupCurrent = &group
		}

		if groupOld == nil {
			return slog.Group(
				"subject",
				slog.String("objecttype", subjectType.String()),
				slog.Group(
					"current",
					*groupCurrent,
				),
			)
		} else {
			return slog.Group(
				"subject",
				slog.String("objecttype", subjectType.String()),
				slog.Group(
					"current",
					*groupCurrent,
				),
				slog.Group(
					"old",
					*groupOld,
				),
			)
		}
	}

	var actorData slog.Attr
	if actor != nil {
		actorData = getActorData(actor.ID, actor.Username, actor.IsAdmin)
	} else {
		actorData = getActorData("nil", "nil", false)
	}
	LogAuditEvent(
		true,
		targetPath,
		srcIp,
		eventType.String(),
		actorData,
		getSubject(subjectCurrent, subjectOld),
	)
}