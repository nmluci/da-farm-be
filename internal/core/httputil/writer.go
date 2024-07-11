package httputil

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/da-farm-be/internal/core/errs"
	"github.com/nmluci/da-farm-be/internal/core/httpres"
)

// WriteSuccessResponse serialized response data
func WriteSuccessResponse(ec echo.Context, status int, data any) error {
	return ec.JSON(status, httpres.BaseResponse{
		Data:   data,
		Errors: nil,
	})
}

// WriteErrorResponse serialized Go's error into standardized error code
func WriteErrorResponse(ec echo.Context, err error) error {
	res := errs.GetErrorResp(err)

	return ec.JSON(res.Status, httpres.BaseResponse{
		Data: nil,
		Errors: &httpres.ErrorResponse{
			Code: res.Code,
			Msg:  res.Msg,
		},
	})
}

// WriteErrorResponseWithStatus serialized Go's error into standardized error code with custom status
func WriteErrorResponseWithStatus(ec echo.Context, status int, err error) error {
	res := errs.GetErrorResp(err)

	return ec.JSON(status, httpres.BaseResponse{
		Data: nil,
		Errors: &httpres.ErrorResponse{
			Code: res.Code,
			Msg:  res.Msg,
		},
	})
}
