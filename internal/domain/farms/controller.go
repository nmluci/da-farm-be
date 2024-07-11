package farms

import "github.com/labstack/echo/v4"

type FarmController struct {
	svc FarmService
}

func NewController(svc FarmService) *FarmController {
	return &FarmController{
		svc: svc,
	}
}

func (fc *FarmController) Route(grp *echo.Group) {
	subrouter := grp.Group("/farms")

	subrouter.GET("", HandleGetAllFarm(fc.svc.GetAll))
	subrouter.OPTIONS("", HandleGetAllFarm(fc.svc.GetAll))

	return
}
