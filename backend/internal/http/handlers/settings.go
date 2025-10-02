package handlers

import (
	"net/http"

	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Settings struct{ svc *service.Services }

func NewSettings(s *service.Services) *Settings { return &Settings{svc: s} }

func (h *Settings) RegisterProtected(v1 *echo.Group) {
	v1.GET("/settings", h.List)
	v1.PATCH("/settings", h.Patch)
}

func (h *Settings) List(c echo.Context) error {
	ctx := c.Request().Context()
	out, err := h.svc.Settings.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list settings"})
	}
	return c.JSON(http.StatusOK, out)
}

type patchRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (h *Settings) Patch(c echo.Context) error {
	var req patchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	ctx := c.Request().Context()
	if req.Key == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "key required"})
	}
	if err := h.svc.Settings.Set(ctx, req.Key, req.Value); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
