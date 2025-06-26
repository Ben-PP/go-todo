package routes

import (
	"go-todo/controllers"
	"go-todo/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	userController *controllers.UserController
}

func NewRouteUser(userController *controllers.UserController) *UserRoutes {
	return &UserRoutes{userController}
}

func (ur *UserRoutes) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/user")
	router.POST("/", ur.userController.CreateUser)
	router.PATCH("/:id", middleware.JwtAuthMiddleware(), ur.userController.UpdateUser)
	router.DELETE("/:id", middleware.JwtAuthMiddleware(), ur.userController.DeleteUser)
}