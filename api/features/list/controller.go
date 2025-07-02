package list

import (
	"context"
	db "go-todo/db/sqlc"
)

type ListController struct {
	db *db.Queries
	ctx context.Context
}

func NewController(db *db.Queries, ctx context.Context) *ListController {
	return &ListController{db: db, ctx: ctx}
}