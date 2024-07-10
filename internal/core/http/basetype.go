package http

type ErrorResponse struct {
	Status int    `json:"-"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
}

type BaseResponse struct {
	Data   any            `json:"data"`
	Errors *ErrorResponse `json:"error"`
}
