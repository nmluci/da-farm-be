package domain

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	ecMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nmluci/da-farm-be/internal/core/middleware"
	"github.com/nmluci/da-farm-be/internal/domain/farms"
	"github.com/nmluci/da-farm-be/internal/domain/ping"
	"github.com/nmluci/da-farm-be/internal/domain/ponds"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitDomain(logger zerolog.Logger, db *sqlx.DB, ec *echo.Echo) {
	// initialize swagger api route
	ec.GET("/api/swagger/*", echoSwagger.WrapHandler)

	// initialize root for backend API
	root := ec.Group("/api/v1",
		ecMiddleware.RequestIDWithConfig(ecMiddleware.RequestIDConfig{Generator: uuid.NewString}),
		middleware.RequestBodyLogger(&logger),
		middleware.RequestLogger(&logger),
		middleware.HandlerLogger(&logger),
		ecMiddleware.CORS(),
	)

	// repository
	farmRepository := farms.NewRepository(db)
	pondRepository := ponds.NewRepository(db)

	// services
	pingService := ping.NewService()
	farmService := farms.NewService(farmRepository)
	pondService := ponds.NewService(pondRepository)

	// handler
	ping.NewController(pingService).Route(root)
	farms.NewController(farmService).Route(root)
	ponds.NewController(pondService).Route(root)
}
