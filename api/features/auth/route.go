package auth

import (
	"go-todo/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authController *AuthController
}

func NewRoutes(authController *AuthController) *AuthRoutes {
	return &AuthRoutes{authController}
}

func (ar *AuthRoutes) Register(rg *gin.RouterGroup) {
	router := rg.Group("/auth")
	router.POST("/login", ar.authController.Login)
	router.POST("/logout", middleware.JwtAuthMiddleware(), ar.authController.Logout)
	router.POST("/refresh", ar.authController.Refresh)
	router.POST("/update-password",middleware.JwtAuthMiddleware(), ar.authController.UpdatePassword)
}