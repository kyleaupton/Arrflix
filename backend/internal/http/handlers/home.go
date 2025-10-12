package handlers

import (
	"net/http"

	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Rail struct{ svc *service.Services }

func NewRails(s *service.Services) *Rail { return &Rail{svc: s} }

func (h *Rail) RegisterProtected(v1 *echo.Group) {
	v1.GET("/home", h.GetHome)
}

// GetRails
// @Summary Get rails
// @Tags    home
// @Produce json
// @Success 200 {array} model.Rail
// @Router  /v1/home [get]
func (h *Rail) GetHome(c echo.Context) error {
	ctx := c.Request().Context()
	rails, err := h.svc.Rails.GetRails(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get rails"})
	}
	return c.JSON(http.StatusOK, rails)
}
