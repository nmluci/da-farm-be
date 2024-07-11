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
//	@Success	200	{object}	ListFarmResponse
//	@Failure	404	{object}	httpres.ErrorResponse
//	@Failure	500	{object}	httpres.ErrorResponse
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
