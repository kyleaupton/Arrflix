package handlers

import (
	"net/http"

	"github.com/kyleaupton/Arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Settings struct{ svc *service.Services }

func NewSettings(s *service.Services) *Settings { return &Settings{svc: s} }

func (h *Settings) RegisterProtected(v1 *echo.Group) {
	v1.GET("/settings", h.List)
	v1.PATCH("/settings", h.Patch)
}

// SettingsListResponse represents a map of setting keys to their values.
// Values may be string, bool, number, or object depending on the registered type.
type SettingsListResponse map[string]any

// List returns all settings with defaults applied.
//
// @Summary List settings
// @Tags    settings
// @Produce json
// @Success 200 {object} handlers.SettingsListResponse
// @Router  /v1/settings [get]
func (h *Settings) List(c echo.Context) error {
	ctx := c.Request().Context()
	out, err := h.svc.Settings.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list settings"})
	}
	return c.JSON(http.StatusOK, out)
}

// PatchRequest updates a single setting value.
// The type of value must match the server-side registry for the given key.
type PatchRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Patch updates a single setting value.
//
// @Summary Update a setting
// @Tags    settings
// @Accept  json
// @Param   payload body  handlers.PatchRequest true "Patch request"
// @Success 204  {string} string ""
// @Failure 400  {object} map[string]string
// @Router  /v1/settings [patch]
func (h *Settings) Patch(c echo.Context) error {
	var req PatchRequest
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
