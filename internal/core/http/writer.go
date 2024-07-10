package http

import (
	"github.com/labstack/echo/v4"
)

// WriteSuccessResponse serialized response data
func WriteSuccessResponse(ec echo.Context, status int, data any) error {
	return ec.JSON(status, BaseResponse{
		Data:   data,
		Errors: nil,
	})
}

// WriteErrorResponse serialized Go's error into standardized error code
func WriteErrorResponse(ec echo.Context, status int, err error) error {
	return ec.JSON(status, BaseResponse{
		Data: nil,
		Errors: &ErrorResponse{
			Code: "DA-01",
			Msg:  err.Error(),
		},
	})
}
