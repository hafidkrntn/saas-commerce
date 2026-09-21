package pagination

// Resolved contains computed pagination values ready for use in queries.
type Resolved struct {
	Page   int
	Limit  int
	Offset int
}

// Pagination is the metadata returned in paginated responses.
type Pagination struct {
	Total     int `json:"total"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
}

// for count data
type CountDetail struct {
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Value float64 `json:"value"`
}

// Response is a generic paginated response container.
type Response[T any] struct {
	Data       []T           `json:"data"`
	Count      []CountDetail `json:"count"`
	Pagination Pagination    `json:"pagination"`
}

// Resolve computes valid page, limit, and offset values from raw input.
func Resolve(page, limit int) Resolved {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return Resolved{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

// NewResponse creates a paginated response from data and total count.
func NewResponse[T any](data []T, total, page, limit int) Response[T] {
	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	if data == nil {
		data = []T{}
	}

	return Response[T]{
		Data: data,
		Pagination: Pagination{
			Total:     total,
			Page:      page,
			Limit:     limit,
			TotalPage: totalPages,
		},
	}
}
