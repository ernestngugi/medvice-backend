package router

import (
	"net/http"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/repository"
	"github.com/ernestngugi/medvice-backend/internal/services"
	"github.com/ernestngugi/medvice-backend/internal/web/api/todo"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
) *AppRouter {

	router := gin.Default()

	appRouter := router.Group("/v1")

	todoRepository := repository.NewTodoRepository()

	todoService := services.NewTodoService(todoRepository)

	todo.AddOpenEndpoints(appRouter, dB, todoService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{
		router,
	}
}
