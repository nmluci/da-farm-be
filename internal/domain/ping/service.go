package ping

import (
	"context"

	"github.com/rs/zerolog"
)

type PingService interface {
	Ping(context.Context) string
}

type pingService struct{}

func NewService() PingService {
	return &pingService{}
}

func (svc *pingService) Ping(ctx context.Context) string {
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("ping recv'd")

	return "hello world"
}
