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

func BuildTodo() *Todo {
	return &Todo{
		Title:       faker.RandomString(5),
		Description: faker.Lorem().Paragraph(2),
	}
}
