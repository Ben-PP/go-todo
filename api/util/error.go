package util

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