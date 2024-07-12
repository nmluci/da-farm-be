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

const (
	farmBasepath = "/farms"
	farmIDPath   = "/:farmID"
)

func (fc *FarmController) Route(grp *echo.Group) {
	subrouter := grp.Group(farmBasepath)

	subrouter.GET("", HandleGetAllFarm(fc.svc.GetAll))
	subrouter.OPTIONS("", HandleGetAllFarm(fc.svc.GetAll))
	subrouter.GET(farmIDPath, HandleGetOneFarm(fc.svc.GetOne))
	subrouter.OPTIONS(farmIDPath, HandleGetOneFarm(fc.svc.GetOne))
	subrouter.POST("", HandleCreateFarm(fc.svc.Create))
	subrouter.OPTIONS("", HandleCreateFarm(fc.svc.Create))
	subrouter.PUT(farmIDPath, HandleUpdateFarm(fc.svc.Update))
	subrouter.OPTIONS(farmIDPath, HandleUpdateFarm(fc.svc.Update))
	subrouter.DELETE(farmIDPath, HandleDeleteFarm(fc.svc.Delete))
	subrouter.OPTIONS(farmIDPath, HandleDeleteFarm(fc.svc.Delete))

	return
}
