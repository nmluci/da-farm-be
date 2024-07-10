// server package act as entry point for backend service
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/da-farm-be/docs"
	_ "github.com/nmluci/da-farm-be/docs"
	"github.com/nmluci/da-farm-be/internal/config"
	"github.com/nmluci/da-farm-be/internal/database/postgres"
	"github.com/nmluci/da-farm-be/internal/domain"
	"github.com/nmluci/da-farm-be/internal/logger"
)

// @title			DA Farm Backend
// @version		1.0
// @description	Simple API to manage Farms and Ponds
// @termsOfService	http://swagger.io/terms/
//
// @BasePath		/api/v1
func main() {
	// bootstrapping
	config := config.New()
	logger := logger.New(config)

	docs.SwaggerInfo.Host = config.ServiceAddress

	db, err := postgres.New(logger, config.PostgresConf)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialized db")
	}

	ec := echo.New()
	ec.HideBanner = true
	ec.HidePort = true

	domain.InitDomain(logger, db, ec)

	logger.Info().Msgf("starting service, listening at %s", config.ServiceAddress)
	if err := ec.Start(config.ServiceAddress); err != nil {
		logger.Error().Err(err).Msg("failed to start service")
	}
}
