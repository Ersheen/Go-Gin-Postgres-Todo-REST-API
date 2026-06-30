package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	// "honnef.co/go/tools/config"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Todo api is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/auth/register", handlers.TodoUsersHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protected := router.Group("/todo")
	protected.Use(middleware.AuthMiddleware(cfg))

	protected.POST("/", handlers.CreateTodoHandler(pool))
	protected.GET("/", handlers.GetAllTodosHandler(pool))
	protected.GET("/:id", handlers.GetTodoHandler(pool))
	protected.PATCH("/:id", handlers.UpdateTodoHandler(pool))
	protected.DELETE("/:id", handlers.DeleteTodoHandler(pool))

	// router.GET("/querytodo", handlers.GetTodoByQueryHandler(pool))
	router.Run(":" + cfg.Port)
}
