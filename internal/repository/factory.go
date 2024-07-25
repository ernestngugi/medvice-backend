package repository

import (
	"context"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/entities"
)

func CreateTodo(ctx context.Context, dB db.DB) (*entities.Todo, error) {
	todo := entities.BuildTodo()
	err := NewTodoRepository().Save(ctx, dB, todo)
	return todo, err
}
