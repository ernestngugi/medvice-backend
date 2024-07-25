package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/entities"
	"github.com/ernestngugi/medvice-backend/internal/forms"
	"github.com/ernestngugi/medvice-backend/internal/providers"
	"github.com/ernestngugi/medvice-backend/internal/repository"
	"github.com/ernestngugi/medvice-backend/internal/utils"
)

const (
	todoKeyPrefix = "medv-api:todo-key:%v"
)

type (
	TodoService interface {
		CompleteTodo(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error)
		CreateTodo(ctx context.Context, dB db.DB, form *forms.CreateTodoForm) (*entities.Todo, error)
		DeleteTodo(ctx context.Context, dB db.DB, todoID int64) error
		TodoByID(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error)
		Todos(ctx context.Context, dB db.DB, filter *forms.Filter) (*entities.TodoList, error)
		UpdateTodo(ctx context.Context, dB db.DB, todoID int64, form *forms.UpdateTodoForm) (*entities.Todo, error)
	}

	todoService struct {
		cacheService   CacheService
		todoRepository repository.TodoRepository
	}
)

func NewTestTodoService(
	redisProvider providers.Redis,
) *todoService {
	cacheService := NewTestCacheService(redisProvider)
	return &todoService{
		cacheService:   cacheService,
		todoRepository: repository.NewTodoRepository(),
	}
}

func NewTodoService(
	cacheService CacheService,
	todoRepository repository.TodoRepository,
) TodoService {
	return &todoService{
		cacheService:   cacheService,
		todoRepository: todoRepository,
	}
}

func (s *todoService) TodoByID(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {

	exist, err := s.cacheService.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return &entities.Todo{}, err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheService.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return &entities.Todo{}, err
		}

		return todo, nil
	}

	todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoService) CreateTodo(ctx context.Context, dB db.DB, form *forms.CreateTodoForm) (*entities.Todo, error) {

	err := utils.ValidateSingleName(form.Title)
	if err != nil {
		return &entities.Todo{}, err
	}

	todo := &entities.Todo{
		Title: form.Title,
	}

	if strings.TrimSpace(form.Description) != "" {
		todo.Description = form.Description
	}

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	err = s.cacheTodo(todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoService) UpdateTodo(ctx context.Context, dB db.DB, todoID int64, form *forms.UpdateTodoForm) (*entities.Todo, error) {

	exist, err := s.cacheService.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return &entities.Todo{}, err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheService.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return &entities.Todo{}, err
		}
	} else {

		todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
		if err != nil {
			return &entities.Todo{}, err
		}

	}

	if form.Title != nil {
		err := utils.ValidateSingleName(*form.Title)
		if err != nil {
			return &entities.Todo{}, err
		}
		todo.Title = *form.Title
	}

	if form.Description != nil {
		if strings.TrimSpace(*form.Description) != "" {
			todo.Description = *form.Description
		}
	}

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	err = s.removeFromCache(todo.ID)
	if err != nil {
		return &entities.Todo{}, err
	}

	err = s.cacheTodo(todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoService) CompleteTodo(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {

	exist, err := s.cacheService.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return &entities.Todo{}, err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheService.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return &entities.Todo{}, err
		}
	} else {

		todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
		if err != nil {
			return &entities.Todo{}, err
		}

	}

	if todo.Completed {
		return &entities.Todo{}, fmt.Errorf("todo has been marked as complete")
	}

	timeNow := time.Now()

	todo.Completed = true
	todo.CompletedAt = &timeNow

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	err = s.removeFromCache(todo.ID)
	if err != nil {
		return &entities.Todo{}, err
	}

	err = s.cacheTodo(todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoService) DeleteTodo(ctx context.Context, dB db.DB, todoID int64) error {

	exist, err := s.cacheService.Exists(s.generateCacheKey(todoID))
	if err != nil {
		return err
	}

	var todo *entities.Todo

	if exist {

		err = s.cacheService.GetCachedValue(s.generateCacheKey(todoID), &todo)
		if err != nil {
			return err
		}
	} else {

		todo, err = s.todoRepository.TodoByID(ctx, dB, todoID)
		if err != nil {
			return err
		}

	}

	if todo.Completed {
		return fmt.Errorf("cannot a todo that has been completed")
	}

	err = s.removeFromCache(todo.ID)
	if err != nil {
		return err
	}

	return s.todoRepository.DeleteTodo(ctx, dB, todo.ID)
}

func (s *todoService) Todos(ctx context.Context, dB db.DB, filter *forms.Filter) (*entities.TodoList, error) {

	todos, err := s.todoRepository.Todos(ctx, dB, filter)
	if err != nil {
		return &entities.TodoList{}, err
	}

	count, err := s.todoRepository.NumberOfTodos(ctx, dB, filter)
	if err != nil {
		return &entities.TodoList{}, err
	}

	todoList := &entities.TodoList{
		Todos:      todos,
		Pagination: entities.NewPagination(count, filter.Page, filter.Per),
	}

	return todoList, nil
}

func (s *todoService) generateCacheKey(todoID int64) string {
	return fmt.Sprintf(todoKeyPrefix, todoID)
}

func (s *todoService) cacheTodo(todo *entities.Todo) error {
	return s.cacheService.CacheValue(s.generateCacheKey(todo.ID), todo)
}

func (s *todoService) removeFromCache(todoID int64) error {
	return s.cacheService.RemoveFromCache(s.generateCacheKey(todoID))
}
