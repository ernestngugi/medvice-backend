package entities

type Pagination struct {
	Count    int  `json:"count"`
	NextPage *int `json:"next_page"`
	NumPages int  `json:"num_pages"`
	Page     int  `json:"page"`
	Per      int  `json:"per"`
	PrevPage *int `json:"prev_page"`
}

func NewPagination(count, page, per int) *Pagination {
	var prevPage, nextPage *int

	if page > 1 {
		pg := page - 1
		prevPage = &pg
	}

	if per < 1 {
		per = 10
	}

	numPages := count / per
	if count == 0 {
		numPages = 1
	} else if count%per != 0 {
		numPages++
	}

	if page < numPages {
		pg := page + 1
		nextPage = &pg
	}

	return &Pagination{
		Count:    count,
		NextPage: nextPage,
		NumPages: numPages,
		Page:     page,
		Per:      per,
		PrevPage: prevPage,
	}
}
