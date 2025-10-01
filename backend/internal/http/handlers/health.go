package handlers

import "github.com/labstack/echo/v4"

type Health struct{}

func NewHealth() *Health { return &Health{} }

func (h *Health) Health(c echo.Context) error {
	return c.String(200, "ok")
}
