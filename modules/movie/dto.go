package movie

type RequestInterface interface {
	Validate() interface{}
}

type Request struct {
	Page    string `form:"page" validate:"blacklist"`
	Limit   int    `form:"limit"`
	Search  string `form:"search"`
	OrderBy string `form:"orderBy"`
	SortBy  string `form:"sortBy"`
}

type Response struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
	Data    any    `json:"data" binding:"required"`
}
