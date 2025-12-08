package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Downloaders struct{ svc *service.Services }

func NewDownloaders(s *service.Services) *Downloaders { return &Downloaders{svc: s} }

func (h *Downloaders) RegisterProtected(v1 *echo.Group) {
	v1.GET("/downloaders", h.List)
	v1.POST("/downloaders", h.Create)
	v1.GET("/downloaders/default/:protocol", h.GetDefault)
	v1.GET("/downloaders/:id", h.Get)
	v1.PUT("/downloaders/:id", h.Update)
	v1.DELETE("/downloaders/:id", h.Delete)
	v1.POST("/downloaders/:id/test", h.Test)
}

// DownloaderCreateRequest payload
type DownloaderCreateRequest struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Protocol   string                 `json:"protocol"`
	URL        string                 `json:"url"`
	Username   *string                `json:"username"`
	Password   *string                `json:"password"`
	ConfigJSON map[string]interface{} `json:"config_json"`
	Enabled    bool                   `json:"enabled"`
	Default    bool                   `json:"default"`
}

// DownloaderUpdateRequest payload
type DownloaderUpdateRequest struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Protocol   string                 `json:"protocol"`
	URL        string                 `json:"url"`
	Username   *string                `json:"username"`
	Password   *string                `json:"password"`
	ConfigJSON map[string]interface{} `json:"config_json"`
	Enabled    bool                   `json:"enabled"`
	Default    bool                   `json:"default"`
}

// List downloaders
// @Summary List downloaders
// @Tags    downloaders
// @Produce json
// @Success 200 {array} dbgen.Downloader
// @Router  /v1/downloaders [get]
func (h *Downloaders) List(c echo.Context) error {
	ctx := c.Request().Context()
	downloaders, err := h.svc.Downloaders.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list downloaders"})
	}
	return c.JSON(http.StatusOK, downloaders)
}

// Create downloader
// @Summary Create downloader
// @Tags    downloaders
// @Accept  json
// @Produce json
// @Param   payload body handlers.DownloaderCreateRequest true "Create downloader"
// @Success 201 {object} dbgen.Downloader
// @Failure 400 {object} map[string]string
// @Router  /v1/downloaders [post]
func (h *Downloaders) Create(c echo.Context) error {
	var req DownloaderCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	ctx := c.Request().Context()
	downloader, err := h.svc.Downloaders.Create(ctx, req.Name, req.Type, req.Protocol, req.URL, req.Username, req.Password, req.ConfigJSON, req.Enabled, req.Default)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, downloader)
}

// Get downloader
// @Summary Get downloader
// @Tags    downloaders
// @Produce json
// @Success 200 {object} dbgen.Downloader
// @Router  /v1/downloaders/{id} [get]
func (h *Downloaders) Get(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	downloader, err := h.svc.Downloaders.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, downloader)
}

// Update downloader
// @Summary Update downloader
// @Tags    downloaders
// @Accept  json
// @Produce json
// @Param   id path string true "Downloader ID"
// @Param   payload body handlers.DownloaderUpdateRequest true "Update downloader"
// @Success 200 {object} dbgen.Downloader
// @Failure 400 {object} map[string]string
// @Router  /v1/downloaders/{id} [put]
func (h *Downloaders) Update(c echo.Context) error {
	var req DownloaderUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	downloader, err := h.svc.Downloaders.Update(ctx, id, req.Name, req.Type, req.Protocol, req.URL, req.Username, req.Password, req.ConfigJSON, req.Enabled, req.Default)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, downloader)
}

// Delete downloader
// @Summary Delete downloader
// @Tags    downloaders
// @Param   id path string true "Downloader ID"
// @Success 204 {string} string ""
// @Router  /v1/downloaders/{id} [delete]
func (h *Downloaders) Delete(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	if err := h.svc.Downloaders.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete"})
	}
	return c.NoContent(http.StatusNoContent)
}

// GetDefault downloader by protocol
// @Summary Get default downloader by protocol
// @Tags    downloaders
// @Produce json
// @Success 200 {object} dbgen.Downloader
// @Router  /v1/downloaders/default/{protocol} [get]
func (h *Downloaders) GetDefault(c echo.Context) error {
	protocol := c.Param("protocol")
	if protocol != "torrent" && protocol != "usenet" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "protocol must be 'torrent' or 'usenet'"})
	}
	ctx := c.Request().Context()
	downloader, err := h.svc.Downloaders.GetDefault(ctx, protocol)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, downloader)
}

// Test downloader connection
// @Summary Test downloader connection
// @Tags    downloaders
// @Param   id path string true "Downloader ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router  /v1/downloaders/{id}/test [post]
func (h *Downloaders) Test(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	downloader, err := h.svc.Downloaders.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	// Test connection by attempting to get version
	switch downloader.Type {
	case "qbittorrent":
		// For now, just return success - full test can be implemented later
		return c.JSON(http.StatusOK, map[string]string{"status": "ok", "message": "Connection test not yet implemented"})
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported downloader type"})
	}
}

