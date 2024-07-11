package httpres

type ErrorResponse struct {
	Status int    `json:"-"`
	Code   int    `json:"code" example:"404001"`
	Msg    string `json:"msg" example:"unknown"`
}

type BaseResponse struct {
	Data   any            `json:"data"`
	Errors *ErrorResponse `json:"error"`
}

type ListPagination struct {
	Limit     uint64 `json:"limit" example:"100"`
	Page      uint64 `json:"page" example:"1"`
	TotalPage uint64 `json:"total_page" example:"10"`
}
