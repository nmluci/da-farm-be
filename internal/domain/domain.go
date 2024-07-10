package domain

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	ecMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/nmluci/da-farm-be/internal/core/middleware"
	"github.com/nmluci/da-farm-be/internal/domain/ping"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitDomain(logger zerolog.Logger, db *sqlx.DB, ec *echo.Echo) {
	// initialize swagger api route
	ec.GET("/api/swagger/*", echoSwagger.WrapHandler)

	// initialize root for backend API
	root := ec.Group("/api/v1",
		ecMiddleware.CORS(), ecMiddleware.CSRF(),
		middleware.RequestBodyLogger(&logger),
		middleware.RequestLogger(&logger),
		middleware.HandlerLogger(&logger),
	)

	// services
	pingService := ping.NewService(logger)

	// handler
	ping.NewController(pingService).Route(root)

}
