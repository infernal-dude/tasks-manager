package main

import (
	"fmt"
	"log"

	"tasks-manager/internal/config"
	"tasks-manager/internal/database"
	"tasks-manager/internal/handler"
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
	if err = db.Ping(); err != nil {
		log.Println("Problem in getting answer from database")
	} else {
		fmt.Println("Krasava x2")
	}

	// err = repo.Delete(1)
	// if err != nil {
	// 	fmt.Println("Error in deleting row")
	// }

	//Обновление записи в базе данных
	// task := domain.Task{ID: 1, Title: "Играть еще больше в Helldivers 2", Description: "Недостаточно играешь", Completed: false}
	// err = repo.Update(&task)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	//Добавление нового таска в базу данных
	// task := &domain.Task{
	// 	Title:       "Играть",
	// 	Description: "Helldivers 2",
	// }
	// repo.Create(task)
	// fmt.Println(task)

	//Получение таска по id
	// task, err := repo.GetById(1)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		fmt.Println("No such rows")
	// 	} else {
	// 		fmt.Printf("Error in getting answer: %s\n", err.Error())
	// 	}
	// }
	// fmt.Println(task)

	//Получение среза всех тасков
	// tasks, err := repo.GetAll()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// for i := 0; i < len(tasks); i++ {
	// 	fmt.Println(tasks[i])
	// }

	r := gin.Default()

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	r.POST("/tasks", handler.Create)
	r.GET("/tasks", handler.GetAll)
	r.GET("/tasks/:id", handler.GetById)
	r.PUT("/tasks/:id", handler.Update)
	r.DELETE("/tasks/:id", handler.Delete)

	r.Run(":8080")
}
