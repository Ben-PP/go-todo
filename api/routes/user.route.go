package routes

import (
	"go-todo/controllers"

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
}