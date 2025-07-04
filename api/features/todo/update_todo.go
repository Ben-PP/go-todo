package todo

import (
	"go-todo/schemas"
	"go-todo/util/mycontext"

	"github.com/gin-gonic/gin"
)

func (controller *TodoController) UpdateTodo(ctx *gin.Context) {
	payload := &schemas.UpdateTodo{}

	if ok := mycontext.ShouldBindBodyWithJSON(&payload, ctx); !ok {
		return
	} else if payload.Title == nil &&
		payload.Description == nil &&
		payload.CompleteBefore == nil {
	}

}
