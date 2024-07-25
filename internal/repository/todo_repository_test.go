package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/entities"
	"github.com/ernestngugi/medvice-backend/internal/forms"
	"github.com/ernestngugi/medvice-backend/internal/testutils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTodoRepository(t *testing.T) {

	testDB := db.InitDB()
	defer testDB.Close()

	ctx := context.Background()

	Convey("TestTodoRepository", t, testutils.WithTestDB(ctx, testDB, func(ctx context.Context, dB db.DB) {

		todoRepository := NewTodoRepository()

		Convey("can save a todo", func() {

			todo := &entities.Todo{
				Title:       "Test",
				Description: "description",
			}

			err := todoRepository.Save(ctx, dB, todo)
			So(err, ShouldBeNil)

			So(todo.ID, ShouldNotBeNil)
			So(todo.CreatedAt, ShouldNotBeZeroValue)
			So(todo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can get todo by id", func() {

			todo, err := CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			fmt.Printf("todo %+v", todo)

			foundTodo, err := todoRepository.TodoByID(ctx, dB, todo.ID)
			So(err, ShouldBeNil)

			So(foundTodo.ID, ShouldEqual, todo.ID)
			So(foundTodo.Title, ShouldEqual, todo.Title)
			So(foundTodo.Description, ShouldEqual, todo.Description)
			So(foundTodo.Completed, ShouldBeFalse)
			So(foundTodo.CompletedAt, ShouldBeNil)
			So(foundTodo.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can update a todo", func() {

			todo, err := CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			timeNow := time.Now()

			todo.Completed = true
			todo.CompletedAt = &timeNow

			err = todoRepository.Save(ctx, dB, todo)
			So(err, ShouldBeNil)

			foundTodo, err := todoRepository.TodoByID(ctx, dB, todo.ID)
			So(err, ShouldBeNil)

			So(foundTodo.ID, ShouldEqual, todo.ID)
			So(foundTodo.Title, ShouldEqual, todo.Title)
			So(foundTodo.Description, ShouldEqual, todo.Description)
			So(foundTodo.Completed, ShouldBeTrue)
			So(foundTodo.CompletedAt, ShouldNotBeZeroValue)
			So(foundTodo.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can get todo page 1 per 1", func() {

			todo, err := CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			_, err = CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			foundTodos, err := todoRepository.Todos(ctx, dB, &forms.Filter{Page: 1, Per: 1})
			So(err, ShouldBeNil)

			foundTodo1 := foundTodos[0]

			So(foundTodo1.ID, ShouldEqual, todo.ID)
			So(foundTodo1.Title, ShouldEqual, todo.Title)
			So(foundTodo1.Description, ShouldEqual, todo.Description)
			So(foundTodo1.Completed, ShouldBeFalse)
			So(foundTodo1.CompletedAt, ShouldBeZeroValue)
			So(foundTodo1.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo1.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can get todo page 2 per 1", func() {

			_, err := CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			todo, err := CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			foundTodos, err := todoRepository.Todos(ctx, dB, &forms.Filter{Page: 2, Per: 1})
			So(err, ShouldBeNil)

			foundTodo1 := foundTodos[0]

			So(foundTodo1.ID, ShouldEqual, todo.ID)
			So(foundTodo1.Title, ShouldEqual, todo.Title)
			So(foundTodo1.Description, ShouldEqual, todo.Description)
			So(foundTodo1.Completed, ShouldBeFalse)
			So(foundTodo1.CompletedAt, ShouldBeZeroValue)
			So(foundTodo1.CreatedAt, ShouldNotBeZeroValue)
			So(foundTodo1.UpdatedAt, ShouldNotBeZeroValue)
		})

		Convey("can delete a todo", func ()  {
			
			todo, err := CreateTodo(ctx, dB)
			So(err, ShouldBeNil)

			err = todoRepository.DeleteTodo(ctx, dB, todo.ID)
			So(err, ShouldBeNil)

			_, err = todoRepository.TodoByID(ctx, dB, todo.ID)
			So(err, ShouldNotBeNil)

			So(err.Error(), ShouldContainSubstring, "sql: no rows in result set")
		})
	}))
}
