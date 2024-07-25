package forms

type CreateTodoForm struct {
	Description string `json:"description" binding:"required"`
	Title       string `json:"title"`
}

type UpdateTodoForm struct {
	Description *string `json:"description"`
	Title       *string `json:"title"`
}
