package routes

import (
	"go-todo/controllers"

	"github.com/gin-gonic/gin"
)

type StatusRoutes struct {
    statusController *controllers.StatusController
}

func NewRouteStatus(statusController *controllers.StatusController) *StatusRoutes {
    return &StatusRoutes{statusController}
}

func (sr *StatusRoutes) StatusRoute(rg *gin.RouterGroup) {
    router := rg.Group("/status")
    router.GET("/", sr.statusController.GetStatus)
}