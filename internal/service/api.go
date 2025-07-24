package service

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Sasha125588/event_app/internal/config"
	"github.com/Sasha125588/event_app/internal/env"
	"github.com/Sasha125588/event_app/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type App struct {
	Router      *gin.Engine
	DB          *sql.DB
	TaskService *TaskService
}

func NewApp() (*App, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbConfig := config.NewDatabaseConfig()
	db, err := dbConfig.ConnectDB()
	if err != nil {
		return nil, err
	}

	if err := config.CreateTables(db); err != nil {
		return nil, err
	}

	taskRepo := repository.NewTaskRepository(db)
	subTaskRepo := repository.NewSubTaskRepository(db)

	taskService := NewTaskService(taskRepo, subTaskRepo)

	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://task-hub-ruby.vercel.app", "https://task-hub-ruby.vercel.app/*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	router.Use(cors.New(corsConfig))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	app := &App{
		Router:      router,
		DB:          db,
		TaskService: taskService,
	}

	return app, nil
}

func (a *App) Run() {
	port := env.GetEnvString("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(a.Router.Run(":" + port))
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}

func Run() {
	app, err := NewApp()
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Close()

	app.Run()
}
