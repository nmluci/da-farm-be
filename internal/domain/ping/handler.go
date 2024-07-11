package ping

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nmluci/da-farm-be/internal/core/httputil"
)

type PingHandler func(context.Context) string

// Ping godoc
//
// @Summary check server status
//
//	@Tags		misc
//	@Produce	json
//	@Success	200
//	@Router		/misc/ping [get]
func HandlePing(handler PingHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		return httputil.WriteSuccessResponse(c, http.StatusOK, handler(ctx))
	}
}
