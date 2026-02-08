package handlers

import (
	"net/http"

	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Version struct{ svc *service.Services }

func NewVersion(s *service.Services) *Version { return &Version{svc: s} }

func (h *Version) RegisterPublic(v1 *echo.Group) {
	v1.GET("/version", h.GetVersion)
}

// GetVersion returns build info and update status.
//
// @Summary Get version and update information
// @Tags    version
// @Produce json
// @Success 200 {object} service.VersionInfo
// @Router  /v1/version [get]
func (h *Version) GetVersion(c echo.Context) error {
	ctx := c.Request().Context()
	info, err := h.svc.Version.GetVersionInfo(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get version information",
		})
	}
	return c.JSON(http.StatusOK, info)
}
