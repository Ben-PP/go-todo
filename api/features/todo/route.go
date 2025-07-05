package todo

import (
	"go-todo/middleware"

	"github.com/gin-gonic/gin"
)

type TodoRoutes struct {
	todoController *TodoController
}

func NewRoutes(todoController *TodoController) *TodoRoutes {
	return &TodoRoutes{todoController}
}

func (routes *TodoRoutes) Register(rg *gin.RouterGroup) {
	router := rg.Group("/list")

	router.Use(middleware.JwtAuthMiddleware())

	router.POST("/", routes.todoController.CreateList)
	router.PATCH("/:listID", routes.todoController.UpdateList)

	todoRouter := router.Group("/:listID/todo")
	todoRouter.POST("/", routes.todoController.CreateTodo)
	todoRouter.PATCH("/:todoID", routes.todoController.UpdateTodo)
	todoRouter.DELETE("/:todoID", routes.todoController.DeleteTodo)
	// TODO Implement GET for todos
}
