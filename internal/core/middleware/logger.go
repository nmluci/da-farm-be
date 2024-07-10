package middleware

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// HandlerLogger tag every request with request-id for traceability
func HandlerLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			l := logger.With().Logger()
			l.UpdateContext(func(cl zerolog.Context) zerolog.Context {
				return cl.Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID))
			})

			c.SetRequest(c.Request().WithContext(l.WithContext(c.Request().Context())))
			return next(c)
		}

	}
}

// RequestLogger log every request received by backend
func RequestLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Str("latency", v.Latency.String()).
				Str("protocol", v.Protocol).
				Str("remoteIP", v.RemoteIP).
				Str("host", v.Host).
				Str("method", v.Method).
				Str("URI", v.URI).
				Str("route-path", v.RoutePath).
				Str("user-agent", v.UserAgent).
				Int("status", v.Status).Msg("request")
			return nil
		},
		LogLatency:   true,
		LogProtocol:  true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogRoutePath: true,
		LogUserAgent: true,
		LogStatus:    true,
	})
}

// RequestBodyLogger log every request's body received by backend
func RequestBodyLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, in []byte, out []byte) {
			loggerInfo := logger.Info()

			loggerInfo = loggerInfo.Str("request-id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Any("request-header", c.Request().Header)
			if string(in) != "" {
				compactJson := &bytes.Buffer{}
				if err := json.Compact(compactJson, in); err != nil {
					logger.Warn().Err(err).Msg("failed to minify json req")
				} else {
					loggerInfo = loggerInfo.RawJSON("request-body", compactJson.Bytes())
				}
			}

			loggerInfo = loggerInfo.Any("response-header", c.Response().Header())
			if string(out) != "" {
				if !strings.Contains(c.Response().Header().Get("Content-Type"), "application/json") {
					loggerInfo = loggerInfo.Str("response-body", "<non-json response>")
				} else {
					compactJson := &bytes.Buffer{}
					if err := json.Compact(compactJson, out); err != nil {
						logger.Warn().Err(err).Msg("failed to minify json req")
					} else {
						loggerInfo = loggerInfo.RawJSON("response-body", compactJson.Bytes())
					}
				}
			}

			loggerInfo.Send()
		},
	})
}
