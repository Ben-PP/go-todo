package routes

import (
	"go-todo/controllers"
	"go-todo/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authController * controllers.AuthController
}

func NewRouteAuth(authController *controllers.AuthController) *AuthRoutes {
	return &AuthRoutes{authController}
}

func (ar *AuthRoutes) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/login", ar.authController.Login)
	router.POST("/logout", ar.authController.Logout)
	router.POST("/refresh", ar.authController.Refresh)
	router.POST("/update-password",middleware.JwtAuthMiddleware(), ar.authController.UpdatePassword)
}