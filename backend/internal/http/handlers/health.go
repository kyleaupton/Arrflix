package handlers

import "github.com/labstack/echo/v4"

type Health struct{}

func NewHealth() *Health { return &Health{} }

// Register attaches the health route to the router.
func (h *Health) RegisterPublic(e *echo.Echo) {
	e.GET("/health", h.Health)
}

// Health
// @Summary Health check
// @Tags    health
// @Produce text/plain
// @Success 200 {string} string "ok"
// @Router  /health [get]
func (h *Health) Health(c echo.Context) error {
	return c.String(200, "ok")
}
