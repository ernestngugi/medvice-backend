package entities

import (
	"time"

	"syreclabs.com/go/faker"
)

type Todo struct {
	Identifier
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	Timestamps
}

type TodoList struct {
	Todos      []*Todo    `json:"todos"`
	Pagination *Pagination `json:"pagination"`
}

func BuildTodo() *Todo {
	return &Todo{
		Title:       faker.RandomString(5),
		Description: faker.Lorem().Paragraph(2),
	}
}
