package handlers

import "github.com/labstack/echo/v4"

type Health struct{}

func NewHealth() *Health { return &Health{} }

func (h *Health) Health(c echo.Context) error {
	return c.String(200, "ok")
}

// Register attaches the health route to the router.
func (h *Health) Register(e *echo.Echo) {
	e.GET("/health", h.Health)
}
