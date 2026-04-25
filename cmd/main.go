package main

import (
	"log"

	"tasks-manager/internal/config"
	"tasks-manager/internal/database"
	"tasks-manager/internal/handler"
	"tasks-manager/internal/middleware"
	"tasks-manager/internal/repository"
	"tasks-manager/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Problem with getting configs for DataBase")
	}
	db := database.NewPostgres(cfg)
	database.RunMigrations(db.DB)

	r := gin.Default()

	taskRepository := repository.NewRepository(db)
	taskService := service.NewService(taskRepository)
	taskHandler := handler.NewHandler(taskService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	tasks := r.Group("/tasks")
	users := r.Group("/user")
	tasks.Use(middleware.AuthMiddleware())
	users.Use(middleware.AuthMiddleware())
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/user/username/:username", userHandler.GetByUsername)
	r.GET("/user/id/:id", userHandler.GetById)
	r.PUT("/user/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)
	tasks.POST("/", taskHandler.Create)
	tasks.GET("/", taskHandler.GetAll)
	tasks.GET("/:id", taskHandler.GetById)
	tasks.PUT("/:id", taskHandler.Update)
	tasks.DELETE("/:id", taskHandler.Delete)

	r.Run(":8080")
}
