package telemetry

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/da-farm-be/internal/core/httputil"
)

type GetRequestMetricHandler func(context.Context) (*ListRequestMetricResponse, error)

// Get Request Metrics
//
//	@Summary	get request metrics for all registered API
//	@Tags		Misc
//	@Produce	json
//	@Success	200	{object}	ListRequestMetricResponse
//@Failure	404		{object}	httpres.ErrorResponse
//	@Router		/telemetry/request-metrics [get]
func HandleGetRequestMetric(handler GetRequestMetricHandler) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := c.Request().Context()

		data, err := handler(ctx)
		if err != nil {
			return httputil.WriteErrorResponse(c, err)
		}

		return httputil.WriteSuccessResponse(c, http.StatusOK, data)
	}
}
