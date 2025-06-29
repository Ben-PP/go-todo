package main

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"go-todo/controllers"
	db "go-todo/db/sqlc"
	"go-todo/logging"
	"go-todo/middleware"
	"go-todo/routes"
	"go-todo/util"

	"github.com/jackc/pgx/v5"
)

var (
    ctx context.Context
)


func main() {
    appLogger := logging.GetLogger()
    slog.SetDefault(appLogger)

    config, err := util.GetConfig()
    if err != nil {
        _, file, line, _ := runtime.Caller(1)
        logging.LogError(err, fmt.Sprintf("%v: %d", file, line), "Failed to load config.")
        return
    }
    
    conn, err := pgx.Connect(context.Background(), config.DbUrl)
    if err != nil {
        _, file, line, _ := runtime.Caller(1)
        logging.LogError(err, fmt.Sprintf("%v: %d", file, line), "Failed to connect to database.")
        return
    } else {
        fmt.Println("Connected to database")
    }

    defer conn.Close(ctx)

    mydb := db.New(conn)

    authController := controllers.NewAuthController(mydb, ctx)
    authRoutes := routes.NewRouteAuth(authController)
    statusController := controllers.NewStatusController(mydb, ctx)
    statusRoutes := routes.NewRouteStatus(statusController)
    userController := controllers.NewUserController(mydb, ctx)
    userRoutes := routes.NewRouteUser(userController)

    router := gin.Default()

    router.Use(middleware.Logger())
    router.Use(middleware.ErrorHandlerMiddleware())
    router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
            return fmt.Sprintf("%s - [%s] \"%s %s %s %d \"%s\" %s\"\n",
        param.ClientIP,
        param.TimeStamp.Format(time.RFC1123),
        param.Method,
        param.Path,
        param.Request.Proto,
        param.StatusCode,
        param.Request.UserAgent(),
        param.ErrorMessage,
    )
    }))
    
    {
        v1 := router.Group("/api/v1")
        authRoutes.UserRoute(v1)
        statusRoutes.StatusRoute(v1)
        userRoutes.UserRoute(v1)    
    }


    slog.Info("Starting server.")
    router.Run("localhost:8000")
}

