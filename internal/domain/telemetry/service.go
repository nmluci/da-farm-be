package telemetry

import (
	"context"

	"github.com/nmluci/da-farm-be/internal/core/errs"
	"github.com/rs/zerolog"
)

// TelemetryService contains public API available to be interacted with
type RequestMetricService interface {
	GetSummary(context.Context) (*ListRequestMetricResponse, error)
	StoreRequestLog(context.Context, *RequestMetricPayload) error
}

type requestMetricService struct {
	repo RequestMetricRepository
}

// NewService return an instance of TelemetryService
func NewService(repo RequestMetricRepository) RequestMetricService {
	return &requestMetricService{
		repo: repo,
	}
}

func (svc *requestMetricService) GetSummary(ctx context.Context) (res *ListRequestMetricResponse, err error) {
	logger := zerolog.Ctx(ctx)

	reqs, err := svc.repo.GetSummary(ctx)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}

	if len(reqs) == 0 {
		return nil, errs.ErrNotFound
	}

	res = &ListRequestMetricResponse{
		Metrics: []*RequestMetric{},
	}

	for _, req := range reqs {
		res.Metrics = append(res.Metrics, &RequestMetric{
			Endpoint:        req.Endpoint,
			MinLatency:      req.MinLatency,
			AvgLatency:      req.AvgLatency,
			MaxLatency:      req.MaxLatency,
			UniqueUserAgent: req.UniqueUA,
			Count:           req.Count,
		})
	}

	return
}

func (svc *requestMetricService) StoreRequestLog(ctx context.Context, payload *RequestMetricPayload) (err error) {
	return svc.repo.Store(ctx, &RequestLogType{
		Endpoint:    payload.Endpoint,
		Latency:     payload.Latency,
		UserAgent:   payload.UserAgent,
		RequestedAt: payload.RequestedAt,
	})
}
