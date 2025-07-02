package list

import (
	"go-todo/middleware"

	"github.com/gin-gonic/gin"
)

type ListRoutes struct {
	listController *ListController
}

func NewRoutes(listController *ListController) *ListRoutes {
	return &ListRoutes{listController}
}

func (routes *ListRoutes) Register(rg *gin.RouterGroup) {
	router := rg.Group("/list")

	router.Use(middleware.JwtAuthMiddleware())

	router.POST("/", routes.listController.CreateList)
}