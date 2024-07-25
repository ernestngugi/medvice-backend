package forms

type Filter struct {
	Page  int
	Per   int
	Valid *bool
}

func (f *Filter) NoPagination() *Filter {
	return &Filter{
		Valid: f.Valid,
	}
}
