package router

import (
	"net/http"
	"os"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/providers"
	"github.com/ernestngugi/medvice-backend/internal/repository"
	"github.com/ernestngugi/medvice-backend/internal/services"
	"github.com/ernestngugi/medvice-backend/internal/web/api/todo"
	"github.com/ernestngugi/medvice-backend/internal/web/middleware"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
	redisManager providers.Redis,
) *AppRouter {

	if os.Getenv("ENVIRONMENT") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	defaultMiddlewares := middleware.DefaultMiddlewares()
	router.Use(defaultMiddlewares...)

	appRouter := router.Group("/v1")

	todoRepository := repository.NewTodoRepository()

	cacheService := services.NewCacheService(redisManager)

	todoService := services.NewTodoService(
		cacheService,
		todoRepository,
	)

	todo.AddOpenEndpoints(appRouter, dB, todoService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{
		router,
	}
}
