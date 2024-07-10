package ping

import "github.com/labstack/echo/v4"

type PingController struct {
	svc PingService
}

func NewController(svc PingService) *PingController {
	return &PingController{
		svc: svc,
	}
}

func (pc *PingController) Route(grp *echo.Group) {
	subrouter := grp.Group("/misc")

	subrouter.GET("/ping", HandlePing(pc.svc.Ping))
	subrouter.OPTIONS("/ping", HandlePing(pc.svc.Ping))

	return
}
