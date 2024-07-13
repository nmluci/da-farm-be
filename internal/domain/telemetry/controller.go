package telemetry

import "github.com/labstack/echo/v4"

type RequestMetricController struct {
	svc RequestMetricService
}

func NewController(svc RequestMetricService) *RequestMetricController {
	return &RequestMetricController{svc: svc}
}

const (
	requestMetricBasepath = "/telemetry/request-metrics"
)

func (mc *RequestMetricController) Route(grp *echo.Group) {
	subrouter := grp.Group(requestMetricBasepath)

	subrouter.GET("", HandleGetRequestMetric(mc.svc.GetSummary))
	subrouter.OPTIONS("", HandleGetRequestMetric(mc.svc.GetSummary))
}
