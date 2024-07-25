package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/entities"
	"github.com/ernestngugi/medvice-backend/internal/forms"
	"github.com/ernestngugi/medvice-backend/internal/mocks"
	"github.com/ernestngugi/medvice-backend/internal/repository"
	"github.com/ernestngugi/medvice-backend/internal/services"
	"github.com/ernestngugi/medvice-backend/internal/testutils"
	"github.com/ernestngugi/medvice-backend/internal/web/middleware"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"syreclabs.com/go/faker"
)

func TestTodoEndpoints(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	ctx := context.Background()

	redisProvider := mocks.NewMockRedisProvider()

	todoService := services.NewTestTodoService(redisProvider)

	Convey("TestTodoEndpoints", t, testutils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		testRouter := gin.Default()
		testRouter.Use(middleware.DefaultMiddlewares()...)

		routerGroup := testRouter.Group("")

		AddOpenEndpoints(routerGroup, dB, todoService)

		Convey("can get todo by id", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			w, err := testutils.DoRequest(testRouter, http.MethodGet, fmt.Sprintf("/todo/%v", todo.ID), nil)
			So(err, ShouldBeNil)

			So(w.Code, ShouldEqual, http.StatusOK)

			var foundTodo entities.Todo

			data, err := io.ReadAll(w.Body)
			So(err, ShouldBeNil)

			err = json.Unmarshal(data, &foundTodo)
			So(err, ShouldBeNil)

			So(foundTodo.ID, ShouldEqual, todo.ID)
			So(foundTodo.Title, ShouldEqual, todo.Title)
			So(foundTodo.Description, ShouldEqual, todo.Description)
			So(foundTodo.Completed, ShouldBeFalse)
			So(foundTodo.CompletedAt, ShouldBeZeroValue)
			So(foundTodo.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can mark a todo as complete", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			w, err := testutils.DoRequest(testRouter, http.MethodPost, fmt.Sprintf("/todo/%v", todo.ID), nil)
			So(err, ShouldBeNil)

			So(w.Code, ShouldEqual, http.StatusOK)

			var foundTodo entities.Todo

			data, err := io.ReadAll(w.Body)
			So(err, ShouldBeNil)

			err = json.Unmarshal(data, &foundTodo)
			So(err, ShouldBeNil)

			So(foundTodo.ID, ShouldEqual, todo.ID)
			So(foundTodo.Title, ShouldEqual, todo.Title)
			So(foundTodo.Description, ShouldEqual, todo.Description)
			So(foundTodo.Completed, ShouldBeTrue)
			So(foundTodo.CompletedAt, ShouldNotBeZeroValue)
			So(foundTodo.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can delete todo by id", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			w, err := testutils.DoRequest(testRouter, http.MethodDelete, fmt.Sprintf("/todo/%v", todo.ID), nil)
			So(err, ShouldBeNil)

			So(w.Code, ShouldEqual, http.StatusOK)

			_, err = todoService.TodoByID(ctx, dB, todo.ID)
			So(err, ShouldNotBeNil)

			So(err.Error(), ShouldContainSubstring, "sql: no rows in result set")
		})

		Convey("can create a todo", func() {

			form := &forms.CreateTodoForm{
				Title:       "test",
				Description: faker.Lorem().Paragraph(3),
			}

			w, err := testutils.DoRequest(testRouter, http.MethodPost, "/todo", form)
			So(err, ShouldBeNil)

			So(w.Code, ShouldEqual, http.StatusOK)

			var todo entities.Todo

			data, err := io.ReadAll(w.Body)
			So(err, ShouldBeNil)

			err = json.Unmarshal(data, &todo)
			So(err, ShouldBeNil)

			So(todo.ID, ShouldEqual, todo.ID)
			So(todo.Title, ShouldEqual, todo.Title)
			So(todo.Description, ShouldEqual, todo.Description)
			So(todo.Completed, ShouldBeFalse)
			So(todo.CompletedAt, ShouldBeZeroValue)
			So(todo.CreatedAt, ShouldNotBeZeroValue)
			So(todo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can update a todo", func() {

			todo, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			newTitle := "newtitle"

			form := &forms.UpdateTodoForm{
				Title: &newTitle,
			}

			w, err := testutils.DoRequest(testRouter, http.MethodPut, fmt.Sprintf("/todo/%v", todo.ID), form)
			So(err, ShouldBeNil)

			So(w.Code, ShouldEqual, http.StatusOK)

			var updatedTodo entities.Todo

			data, err := io.ReadAll(w.Body)
			So(err, ShouldBeNil)

			err = json.Unmarshal(data, &updatedTodo)
			So(err, ShouldBeNil)

			So(updatedTodo.ID, ShouldEqual, todo.ID)
			So(updatedTodo.Title, ShouldEqual, "newtitle")
			So(updatedTodo.Description, ShouldEqual, todo.Description)
			So(updatedTodo.Completed, ShouldBeFalse)
			So(updatedTodo.CompletedAt, ShouldBeZeroValue)
			So(updatedTodo.CreatedAt, ShouldNotBeZeroValue)
			So(updatedTodo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can get a list of todos", func() {

			todo1, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			todo2, err := repository.CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			w, err := testutils.DoRequest(testRouter, http.MethodGet, "/todos?page=1&per=20", nil)
			So(err, ShouldBeNil)

			So(w.Code, ShouldEqual, http.StatusOK)

			var todoList entities.TodoList

			data, err := io.ReadAll(w.Body)
			So(err, ShouldBeNil)

			err = json.Unmarshal(data, &todoList)
			So(err, ShouldBeNil)

			foundTodo1 := todoList.Todos[0]
			foundTodo2 := todoList.Todos[1]

			So(foundTodo1.ID, ShouldEqual, todo1.ID)
			So(foundTodo1.Title, ShouldEqual, todo1.Title)
			So(foundTodo1.Description, ShouldEqual, todo1.Description)
			So(foundTodo1.Completed, ShouldBeFalse)
			So(foundTodo1.CompletedAt, ShouldBeZeroValue)
			So(foundTodo1.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo1.UpdatedAt, ShouldNotBeZeroValue)

			So(foundTodo2.ID, ShouldEqual, todo2.ID)
			So(foundTodo2.Title, ShouldEqual, todo2.Title)
			So(foundTodo2.Description, ShouldEqual, todo2.Description)
			So(foundTodo2.Completed, ShouldBeFalse)
			So(foundTodo2.CompletedAt, ShouldBeZeroValue)
			So(foundTodo2.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo2.UpdatedAt, ShouldNotBeZeroValue)

			So(todoList.Pagination.Count, ShouldEqual, 2)
		})
	}))
}
