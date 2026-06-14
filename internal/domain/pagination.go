package domain

type Meta struct {
	NextCursor int  `json:"next_cursor"`
	HasCrusor  bool `json:"has_cursor"`
}

type Pagination struct {
	Entities interface{} `json:"entities"`
	Meta     *Meta       `json:"meta"`
}

type PaginationRequest struct {
	Cursor int `form:"cursor,default=0" json:"cursor"`
	Limit  int `form:"limit,default=20" json:"limit"`
}
