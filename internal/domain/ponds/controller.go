package ponds

import "github.com/labstack/echo/v4"

type PondController struct {
	svc PondService
}

func NewController(svc PondService) *PondController {
	return &PondController{
		svc: svc,
	}
}

const (
	pondBasepath = "/farms/:farmID/ponds"
	pondIDPath   = "/:pondID"
)

func (pc *PondController) Route(grp *echo.Group) {
	subrouter := grp.Group(pondBasepath)

	subrouter.GET("", HandleGetAllPond(pc.svc.GetAll))
	subrouter.OPTIONS("", HandleGetAllPond(pc.svc.GetAll))
	subrouter.GET(pondIDPath, HandleGetOnePond(pc.svc.GetOne))
	subrouter.OPTIONS(pondIDPath, HandleGetOnePond(pc.svc.GetOne))
	subrouter.POST("", HandleCreatePond(pc.svc.Create))
	subrouter.OPTIONS("", HandleCreatePond(pc.svc.Create))
	subrouter.PUT(pondIDPath, HandleUpdatePond(pc.svc.Update))
	subrouter.OPTIONS(pondIDPath, HandleUpdatePond(pc.svc.Update))
	subrouter.DELETE(pondIDPath, HandleDeletePond(pc.svc.Delete))
	subrouter.OPTIONS(pondIDPath, HandleDeletePond(pc.svc.Delete))

	return
}
