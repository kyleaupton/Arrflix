package middlewares

import (
	"net/http"
	"strings"

	"github.com/kyleaupton/Arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

// SetupMode middleware enforces setup mode:
// - If NOT initialized: allow /setup/* routes, block everything else (redirect to setup)
// - If initialized: block /setup/* routes (409 Conflict), allow everything else
func SetupMode(services *service.Services) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			ctx := c.Request().Context()

			// Check initialization status
			initialized, err := services.Setup.IsInitialized(ctx)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "setup check failed"})
			}

			// Setup routes: /api/v1/setup/*
			isSetupRoute := strings.HasPrefix(path, "/api/v1/setup/")

			// Health check is always allowed
			if path == "/health" {
				return next(c)
			}

			if !initialized {
				// SETUP MODE: Only allow setup routes
				if !isSetupRoute {
					return c.JSON(http.StatusServiceUnavailable, map[string]string{
						"error":     "setup required",
						"setup_url": "/api/v1/setup/status",
					})
				}
				return next(c)
			} else {
				// NORMAL MODE: Block setup routes
				if isSetupRoute {
					return c.JSON(http.StatusConflict, map[string]string{"error": "system already initialized"})
				}
				return next(c)
			}
		}
	}
}
