package response

type Response struct {
	Code int  `json:"-"`
	Meta Meta `json:"meta,omitzero"`
	Data any  `json:"data,omitempty"`
}

type Meta struct {
	Message         string            `json:"message"`
	Error           string            `json:"error,omitempty"`
	Pagination      Pagination        `json:"pagination,omitzero"`
	ErrorValidation map[string]string `json:"error_validation,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
