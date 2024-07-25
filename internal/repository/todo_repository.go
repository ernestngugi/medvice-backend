package repository

import (
	"context"
	"fmt"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/entities"
	"github.com/ernestngugi/medvice-backend/internal/forms"
)

const (
	countTodoSQL   = "SELECT COUNT(id) FROM todos"
	getTodoByIDSQL = selectTodoSQL + " WHERE id = $1"
	insertTodoSQL  = "INSERT INTO todos (title, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
	selectTodoSQL  = "SELECT id, title, description, completed, completed_at, created_at, updated_at FROM todos"
	updateTodoSQL  = "UPDATE todos SET title = $1, description = $2, completed = $3, completed_at = $4, updated_at = $5 WHERE id = $6"
)

type (
	TodoRepository interface {
		NumberOfTodos(ctx context.Context, operations db.SQLOperations, filter *forms.Filter) (int, error)
		Save(ctx context.Context, operations db.SQLOperations, todo *entities.Todo) error
		TodoByID(ctx context.Context, operations db.SQLOperations, todoID int64) (*entities.Todo, error)
		Todos(ctx context.Context, operations db.SQLOperations, filter *forms.Filter) ([]*entities.Todo, error)
	}

	todoRepository struct{}
)

func NewTodoRepository() TodoRepository {
	return &todoRepository{}
}

func (r *todoRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	todo *entities.Todo,
) error {

	todo.Touch()

	if todo.IsNew() {

		err := operations.QueryRowContext(
			ctx,
			insertTodoSQL,
			todo.Title,
			todo.Description,
			todo.CreatedAt,
			todo.UpdatedAt,
		).Scan(&todo.ID)
		if err != nil {
			return fmt.Errorf("insert todo query error %v", err)
		}

		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateTodoSQL,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.CompletedAt,
		todo.UpdatedAt,
		todo.ID,
	)
	if err != nil {
		return fmt.Errorf("update todo query error %v", err)
	}

	return nil
}

func (r *todoRepository) TodoByID(
	ctx context.Context,
	operations db.SQLOperations,
	todoID int64,
) (*entities.Todo, error) {

	row := operations.QueryRowContext(
		ctx,
		getTodoByIDSQL,
		todoID,
	)

	return r.scanRow(row)
}

func (r *todoRepository) Todos(
	ctx context.Context,
	operations db.SQLOperations,
	filter *forms.Filter,
) ([]*entities.Todo, error) {

	query := selectTodoSQL
	args := []any{}

	if filter.Per > 0 && filter.Page > 0 {
		query += " LIMIT $1 OFFSET $2"
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	rows, err := operations.QueryContext(ctx, query, args...)
	if err != nil {
		return []*entities.Todo{}, fmt.Errorf("todos query error %v", err)
	}

	defer rows.Close()

	todos := make([]*entities.Todo, 0)

	for rows.Next() {
		todo, err := r.scanRow(rows)
		if err != nil {
			return []*entities.Todo{}, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return []*entities.Todo{}, fmt.Errorf("todos rows error %v", err)
	}

	return todos, nil
}

func (r *todoRepository) NumberOfTodos(
	ctx context.Context,
	operations db.SQLOperations,
	filter *forms.Filter,
) (int, error) {

	query := countTodoSQL
	args := []any{}

	var count int

	err := operations.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("number of todos query error %v", err)
	}

	return count, nil
}

func (r *todoRepository) scanRow(
	rowScanner db.RowScanner,
) (*entities.Todo, error) {

	var todo entities.Todo

	err := rowScanner.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CompletedAt,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return &entities.Todo{}, fmt.Errorf("todo scan row error %v", err)
	}

	return &todo, nil
}
