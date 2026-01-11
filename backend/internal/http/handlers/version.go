package handlers

import (
	"net/http"

	"github.com/kyleaupton/Arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Version struct{ svc *service.Services }

func NewVersion(s *service.Services) *Version { return &Version{svc: s} }

func (h *Version) RegisterPublic(v1 *echo.Group) {
	v1.GET("/version", h.GetVersion)
	v1.GET("/update", h.GetUpdate)
}

// GetVersion returns current build information.
//
// @Summary Get version information
// @Tags    version
// @Produce json
// @Success 200 {object} versioninfo.BuildInfo
// @Router  /v1/version [get]
func (h *Version) GetVersion(c echo.Context) error {
	info := h.svc.Version.GetBuildInfo()
	return c.JSON(http.StatusOK, info)
}

// GetUpdate checks for available updates.
//
// @Summary Check for updates
// @Tags    version
// @Produce json
// @Success 200 {object} service.UpdateInfo
// @Router  /v1/update [get]
func (h *Version) GetUpdate(c echo.Context) error {
	ctx := c.Request().Context()
	info, err := h.svc.Version.CheckUpdate(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to check for updates",
		})
	}
	return c.JSON(http.StatusOK, info)
}
