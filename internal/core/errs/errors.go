package errs

import (
	"errors"
	"net/http"

	"github.com/nmluci/da-farm-be/internal/core/httpres"
)

// alias several common error message
var (
	ErrBadRequest               = errors.New("bad request")
	ErrBrokenUserReq            = errors.New("invalid request")
	ErrDuplicatedResources      = errors.New("entity already existed")
	ErrUnknown                  = errors.New("internal server error")
	ErrNotFound                 = errors.New("entity not found")
	ErrMissingRequiredAttribute = errors.New("attribute is missing")
)

// Errcode: AAA-BB-C
// AAA => HTTP STATUS CODE
// BB = 01 Basic, 02 Business Logic
// C = ErrorID
// Ex: 403021 = 403 (Forbidden) - Business Logic - ID 1
const (
	ErrCodeBadRequest               int = 400011
	ErrCodeMissingRequiredAttribute int = 400012
	ErrCodeInvalidCred              int = 401013
	ErrCodeNoAccess                 int = 403014
	ErrCodeNotFound                 int = 404015
	ErrCodeDuplicatedResources      int = 409016
	ErrCodeBrokenUserReq            int = 422017
	ErrCodeUndefined                int = 500011
)

// aliased HTTP status
const (
	ErrStatusUnknown        = http.StatusInternalServerError
	ErrStatusClient         = http.StatusBadRequest
	ErrStatusNotLoggedIn    = http.StatusUnauthorized
	ErrStatusNoAccess       = http.StatusForbidden
	ErrStatusReqBody        = http.StatusUnprocessableEntity
	ErrStatusNotFound       = http.StatusNotFound
	ErrStatusMissingContext = http.StatusPreconditionFailed
)

var errorMap = map[error]httpres.ErrorResponse{
	ErrUnknown:                  errorResponse(ErrStatusUnknown, ErrCodeUndefined, ErrUnknown),
	ErrBadRequest:               errorResponse(ErrStatusClient, ErrCodeBadRequest, ErrBadRequest),
	ErrDuplicatedResources:      errorResponse(ErrStatusClient, ErrCodeDuplicatedResources, ErrDuplicatedResources),
	ErrBrokenUserReq:            errorResponse(ErrStatusReqBody, ErrCodeBrokenUserReq, ErrBrokenUserReq),
	ErrNotFound:                 errorResponse(ErrStatusNotFound, ErrCodeNotFound, ErrNotFound),
	ErrMissingRequiredAttribute: errorResponse(ErrStatusClient, ErrCodeMissingRequiredAttribute, ErrMissingRequiredAttribute),
}

func errorResponse(status int, code int, err error) httpres.ErrorResponse {
	return httpres.ErrorResponse{
		Status: status,
		Code:   code,
		Msg:    err.Error(),
	}
}

// GetErrorResp return ErrorResponse
func GetErrorResp(err error) httpres.ErrorResponse {
	errResponse, ok := errorMap[err]
	if !ok {
		errResponse = errorMap[ErrUnknown]
	}

	return errResponse
}
