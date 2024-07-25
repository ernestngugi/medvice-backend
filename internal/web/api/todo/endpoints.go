package todo

import (
	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func AddOpenEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	todoService services.TodoService,
) {
	r.POST("/todo", createTodo(dB, todoService))
	r.GET("/todos", listTodo(dB, todoService))
	r.GET("/todo/:id", todoByID(dB, todoService))
	r.PUT("/todo/:id", updateTodo(dB, todoService))
	r.POST("/todo/:id", completeTodo(dB, todoService))
	r.DELETE("/todo/:id", deleteTodo(dB, todoService))
}
