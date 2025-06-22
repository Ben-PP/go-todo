package controllers

import (
	"context"
	db "go-todo/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusController struct{
	db *db.Queries
	ctx context.Context
}

// NewUserController creates a new user controller
func NewStatusController(db *db.Queries, ctx context.Context) *StatusController {
	return &StatusController{db: db, ctx: ctx}
}

func (cc *StatusController) GetStatus(ctx *gin.Context) {

    users, err := cc.db.GetAllUsers(ctx)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed to fetch users", "error": err.Error()})
		return
    }
	hasUsers := len(users) != 0

    ctx.JSON(http.StatusOK, gin.H{"status": "ok", "has_users": hasUsers})
}