package util


type ErrInternalMeta struct {
	OrigErrMessage string
	File string
}

func NewErrInternalMeta(file string, origErrMessage string) *ErrInternalMeta {
	return &ErrInternalMeta{
		File: file,
		OrigErrMessage: origErrMessage,
	}
}

type ErrDatabaseMeta struct {
	File string
	QueryDetails string
	RespondWithStatus int
}

func NewErrDatabaseMeta(file string, queryDetails string, respondWithStatus int) *ErrDatabaseMeta {
	return &ErrDatabaseMeta{
		File: file,
		QueryDetails: queryDetails,
		RespondWithStatus: respondWithStatus,
	}
}