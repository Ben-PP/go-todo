package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"go-todo/controllers"
	db "go-todo/db/sqlc"
	"go-todo/routes"
	"go-todo/util"

	//"github.com/golang-jwt/jwt/v5"

	"github.com/jackc/pgx/v5"
)

var (
    ctx context.Context
)


func main() {
    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }
    
    conn, err := pgx.Connect(context.Background(), config.DbUrl)
    if err != nil {
        fmt.Println("Error connecting to database", err)
    }

    defer conn.Close(ctx)

    mydb := db.New(conn)

    statusController := controllers.NewStatusController(mydb, ctx)
    statusRoutes := routes.NewRouteStatus(statusController)
    userController := controllers.NewUserController(mydb, ctx)
    userRoutes := routes.NewRouteUser(userController)

    router := gin.Default()
    {
        v1 := router.Group("/api/v1")
        statusRoutes.StatusRoute(v1)
        userRoutes.UserRoute(v1)    
    }

    router.Run("localhost:8000")
}

