package ponds

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/da-farm-be/internal/core/httputil"
	"github.com/rs/zerolog"
)

type GetAllPondHandler func(context.Context, *PondRequestQuery) (*ListPondResponse, error)

// Get All Pond godoc
//
//	@Summary	get all pond
//	@Tags		Pond
//	@Produce	json
//	@Success	200		{object}	ListPondResponse
//	@Failure	404		{object}	httpres.ErrorResponse
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Param		farmID	path		int		true	"farm ID"
//	@Param		keyword	query		string	false	"Keyword to search"
//	@Param		limit	query		string	false	"number of entity per page"
//	@Param		page	query		string	false	"n-th page"
//	@Router		/farms/{farmID}/ponds [get]
func HandleGetAllPond(handler GetAllPondHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &PondRequestQuery{}

		if err = c.Bind(params); err != nil {
			logger.Err(err).Send()
			return httputil.WriteErrorResponseWithStatus(c, http.StatusBadRequest, err)
		}

		data, err := handler(ctx, params)
		if err != nil {
			return httputil.WriteErrorResponse(c, err)
		}

		return httputil.WriteSuccessResponse(c, http.StatusOK, data)
	}
}

type GetOnePondHandler func(context.Context, *PondRequestQuery) (*PondResponse, error)

// Get One Pond godoc
//
//	@Summary	get specific pond by ID
//	@Tags		Pond
//	@Produce	json
//	@Param		farmID	path		int	true	"Farm ID"
//	@Param		pondID	path		int	true	"Pond ID"
//	@Success	200		{object}	PondResponse
//	@Failure	404		{object}	httpres.ErrorResponse
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms/{farmID}/ponds/{pondID} [get]
func HandleGetOnePond(handler GetOnePondHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &PondRequestQuery{}

		if err = c.Bind(params); err != nil {
			logger.Err(err).Send()
			return httputil.WriteErrorResponseWithStatus(c, http.StatusBadRequest, err)
		}

		data, err := handler(ctx, params)
		if err != nil {
			return httputil.WriteErrorResponse(c, err)
		}

		return httputil.WriteSuccessResponse(c, http.StatusOK, data)
	}
}

type CreatePondHandler func(context.Context, *PondPayload) error

// CreatePond godoc
//
//	@Summary	create a new pond
//	@Tags		Pond
//	@Accept		json
//	@Produce	json
//	@Param		farmID	path		int			true	"Farm ID"
//	@Param		payload	body		PondPayload	true	"pond payload"
//	@Success	201		{object}	string
//	@Failure	409		{object}	httpres.ErrorResponse	"pond with same name already exists"
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms/{farmID}/ponds [post]
func HandleCreatePond(handler CreatePondHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &PondPayload{}

		if err = c.Bind(params); err != nil {
			logger.Err(err).Send()
			return httputil.WriteErrorResponseWithStatus(c, http.StatusBadRequest, err)
		}

		err = handler(ctx, params)
		if err != nil {
			return httputil.WriteErrorResponse(c, err)
		}

		return httputil.WriteSuccessResponse(c, http.StatusCreated, nil)
	}
}

type UpdatePondHandler func(context.Context, *PondPayload) error

// Update Pond godoc
//
//	@Summary	update pond data
//	@Tags		Pond
//	@Accept		json
//	@Produce	json
//	@Param		farmID	path		int			true	"Farm ID"
//	@Param		pondID	path		int			true	"Pond ID"
//	@Param		payload	body		PondPayload	true	"pond payload"
//	@Success	200		{object}	string
//	@Failure	404		{object}	httpres.ErrorResponse	"pond not existed"
//	@Failure	409		{object}	httpres.ErrorResponse	"duplicated pond found"
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms/{farmID}/ponds/{pondID} [put]
func HandleUpdatePond(handler UpdatePondHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()

		logger := zerolog.Ctx(ctx)
		payload := &PondPayload{}

		if err = c.Bind(payload); err != nil {
			logger.Err(err).Send()
			return httputil.WriteErrorResponseWithStatus(c, http.StatusBadRequest, err)
		}

		err = handler(ctx, payload)
		if err != nil {
			return httputil.WriteErrorResponse(c, err)
		}

		return httputil.WriteSuccessResponse(c, http.StatusOK, nil)
	}
}

type DeletePondHandler func(context.Context, *PondRequestQuery) error

// DeletePond godoc
//
//	@Summary	delete specific pond by ID
//	@Tags		Pond
//	@Produce	json
//	@Param		farmID	path		int	true	"Farm ID"
//	@Param		pondID	path		int	true	"Pond ID"
//	@Success	200		{object}	PondResponse
//	@Failure	404		{object}	httpres.ErrorResponse	"pond not existed"
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms/{farmID}/ponds/{pondID} [delete]
func HandleDeletePond(handler DeletePondHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &PondRequestQuery{}

		if err = c.Bind(params); err != nil {
			logger.Err(err).Send()
			return httputil.WriteErrorResponseWithStatus(c, http.StatusBadRequest, err)
		}

		err = handler(ctx, params)
		if err != nil {
			return httputil.WriteErrorResponse(c, err)
		}

		return httputil.WriteSuccessResponse(c, http.StatusOK, nil)
	}
}
