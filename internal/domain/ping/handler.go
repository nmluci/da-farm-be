package ping

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	inhttp "github.com/nmluci/da-farm-be/internal/core/http"
)

type PingHandler func(context.Context) string

// Ping godoc
//
// @Summary check server status
//	@Tags		misc
//	@Produce	json
//	@Success	200
//	@Router		/misc/ping [get]
func HandlePing(handler PingHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		return inhttp.WriteSuccessResponse(c, http.StatusOK, handler(c.Request().Context()))
	}
}
