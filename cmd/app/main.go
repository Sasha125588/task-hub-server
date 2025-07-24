package main

import (
	"log"

	"github.com/Sasha125588/event_app/docs"
	"github.com/Sasha125588/event_app/internal/handlers"
	"github.com/Sasha125588/event_app/internal/service"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Task Hub API
// @version 1.0
// @description This is a task management server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	docs.SwaggerInfo.Title = "Task Hub API"
	docs.SwaggerInfo.Description = "This is a task management server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	app, err := service.NewApp()
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Close()

	taskHandler := handlers.NewTaskHandler(app.TaskService)

	setupRoutes(app.Router, taskHandler)

	app.Router.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))

	app.Run()
}

func setupRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)

			tasks.POST("/:id/subtasks", taskHandler.CreateSubTask)
			tasks.GET("/:id/subtasks", taskHandler.GetSubTasksByTaskID)
		}

		subtasks := v1.Group("/subtasks")
		{
			subtasks.PUT("/:id", taskHandler.UpdateSubTask)
			subtasks.DELETE("/:id", taskHandler.DeleteSubTask)
		}
	}
}
