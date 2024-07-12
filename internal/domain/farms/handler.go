package farms

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/da-farm-be/internal/core/httputil"
	"github.com/rs/zerolog"
)

type GetAllFarmHandler func(context.Context, *FarmRequestQuery) (*ListFarmResponse, error)

// Get All Farm godoc
//
//	@Summary	get all farm
//	@Tags		Farm
//	@Produce	json
//	@Success	200		{object}	ListFarmResponse
//	@Failure	404		{object}	httpres.ErrorResponse
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Param		keyword	query		string	false	"Keyword to search"
//	@Param		limit	query		string	false	"number of entity per page"
//	@Param		page	query		string	false	"n-th page"
//	@Router		/farms [get]
func HandleGetAllFarm(handler GetAllFarmHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &FarmRequestQuery{}

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

type GetOneFarmHandler func(context.Context, *FarmRequestQuery) (*FarmResponse, error)

// Get One Farm godoc
//
//	@Summary	get specific farm by ID
//	@Tags		Farm
//	@Produce	json
//	@Param		farmID	path		int	true	"Farm ID"
//	@Success	200		{object}	FarmResponse
//	@Failure	404		{object}	httpres.ErrorResponse
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms/{farmID} [get]
func HandleGetOneFarm(handler GetOneFarmHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &FarmRequestQuery{}

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

type CreateFarmHandler func(context.Context, *FarmPayload) error

// CreateFarm godoc
//
//	@Summary	create a new farm
//	@Tags		Farm
//	@Accept		json
//	@Produce	json
//	@Param		payload	body		FarmPayload	true	"farm payload"
//	@Success	201		{object}	string
//	@Failure	409		{object}	httpres.ErrorResponse	"farm with same name already exists"
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms [post]
func HandleCreateFarm(handler CreateFarmHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &FarmPayload{}

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

type UpdateFarmHandler func(context.Context, *FarmPayload) error

// Update Farm godoc
//
//	@Summary	update farm data
//	@Tags		Farm
//	@Accept		json
//	@Produce	json
//	@Param		farmID	path		int			true	"Farm ID"
//	@Param		payload	body		FarmPayload	true	"farm payload"
//	@Success	200		{object}	string
//	@Failure	404		{object}	httpres.ErrorResponse "farm not existed"
//	@Failure	409		{object}	httpres.ErrorResponse "duplicated farm found"
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms/{farmID} [put]
func HandleUpdateFarm(handler UpdateFarmHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()

		logger := zerolog.Ctx(ctx)
		payload := &FarmPayload{}

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

type DeleteFarmHandler func(context.Context, *FarmRequestQuery) error

// DeleteFarm godoc
//
//	@Summary	delete specific farm by ID
//	@Tags		Farm
//	@Produce	json
//	@Param		farmID	path		int	true	"Farm ID"
//	@Success	200		{object}	FarmResponse
//	@Failure	404		{object}	httpres.ErrorResponse "farm not existed"
//	@Failure	500		{object}	httpres.ErrorResponse
//	@Router		/farms/{farmID} [delete]
func HandleDeleteFarm(handler DeleteFarmHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)
		params := &FarmRequestQuery{}

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
