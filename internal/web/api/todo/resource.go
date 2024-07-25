package todo

import (
	"net/http"
	"strconv"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/forms"
	"github.com/ernestngugi/medvice-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func createTodo(
	dB db.DB,
	todoService services.TodoService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.CreateTodoForm

		err := c.BindJSON(&form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		todo, err := todoService.CreateTodo(c.Request.Context(), dB, &form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func completeTodo(
	dB db.DB,
	todoService services.TodoService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		todo, err := todoService.CompleteTodo(c.Request.Context(), dB, todoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func updateTodo(
	dB db.DB,
	todoService services.TodoService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.UpdateTodoForm

		err := c.BindJSON(&form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		todo, err := todoService.UpdateTodo(c.Request.Context(), dB, todoID, &form)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func todoByID(
	dB db.DB,
	todoService services.TodoService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		todo, err := todoService.TodoByID(c.Request.Context(), dB, todoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func deleteTodo(
	dB db.DB,
	todoService services.TodoService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		err = todoService.DeleteTodo(c.Request.Context(), dB, todoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func listTodo(
	dB db.DB,
	todoService services.TodoService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todos, err := todoService.Todos(c.Request.Context(), dB, &forms.Filter{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}
