package user

import (
	"go-todo/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userController *UserController
}

func NewRoutes(userController *UserController) *UserRoutes {
	return &UserRoutes{userController}
}

func (ur *UserRoutes) Register(rg *gin.RouterGroup) {
	router := rg.Group("/user")
	router.POST("/", ur.userController.CreateUser)
	router.PATCH("/:id", middleware.JwtAuthMiddleware(), ur.userController.UpdateUser)
	router.DELETE("/:id", middleware.JwtAuthMiddleware(), ur.userController.DeleteUser)
}