package ping

import (
	"context"

	"github.com/rs/zerolog"
)

type PingService interface {
	Ping(context.Context) string
}

type pingService struct {
	logger zerolog.Logger
}

func NewService(logger zerolog.Logger) PingService {
	return &pingService{
		logger: logger,
	}
}

func (svc *pingService) Ping(_ context.Context) string {
	svc.logger.Info().Msg("ping recv'd")

	return "hello world"
}
