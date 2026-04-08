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

	// Для чего этот метод? То что он делает, можно явно сделать здесь и не скрывать ничего.
	// А для него ещё и целый файл выделился, тоже непонятно зачем. Даже целый пакет!
	// По хорошему при создании нового пакета нужно написать комментарий для чего он создан. Что написать для пакета database?
	// К томуже ты в нём уже делаешь Ping(), а тут делает его ещё раз. Так что это лишний раз доказывает бесполезность этого метода.
	db := database.NewPostgres(cfg)
	if err = db.Ping(); err != nil {
		log.Println("Problem in getting answer from database") // После этого лога приложение не остановится, хотя явно должно. Нужен log.Fatal()
	}

	// Зачем эта переменная, db.DB можно явно передать
	sqlDB := db.DB
	database.RunMigrations(sqlDB)

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

	// Про джин я уже говорил вроде, стоит пользоваться стандартными средствами языка на первых парах
	r := gin.Default()

	taskRepository := repository.NewRepository(db)
	taskService := service.NewService(taskRepository)
	taskHandler := handler.NewHandler(taskService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/user/username/:username", userHandler.GetByUsername)
	r.GET("/user/id/:id", userHandler.GetById)
	r.PUT("/user/:id", userHandler.Update)
	r.DELETE("/user/:id", userHandler.Delete)
	tasks := r.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	tasks.POST("/", taskHandler.Create)
	tasks.GET("/", taskHandler.GetAll)
	tasks.GET("/:id", taskHandler.GetById)
	tasks.PUT("/:id", taskHandler.Update)
	tasks.DELETE("/:id", taskHandler.Delete)

	r.Run(":8080")
}
