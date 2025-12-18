package handlers

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type DownloadJobs struct{ svc *service.Services }

func NewDownloadJobs(s *service.Services) *DownloadJobs { return &DownloadJobs{svc: s} }

func (h *DownloadJobs) RegisterProtected(v1 *echo.Group) {
	v1.GET("/download-jobs", h.List)
	v1.GET("/download-jobs/:id", h.Get)
	v1.DELETE("/download-jobs/:id", h.Cancel)

	v1.GET("/movie/:id/download-jobs", h.ListForMovie)
}

// List download jobs
// @Summary List download jobs
// @Tags    download-jobs
// @Produce json
// @Success 200 {array} dbgen.DownloadJob
// @Router  /v1/download-jobs [get]
func (h *DownloadJobs) List(c echo.Context) error {
	ctx := c.Request().Context()
	out, err := h.svc.DownloadJobs.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list"})
	}
	return c.JSON(http.StatusOK, out)
}

// Get download job
// @Summary Get download job
// @Tags    download-jobs
// @Produce json
// @Param   id path string true "Job ID (uuid)"
// @Success 200 {object} dbgen.DownloadJob
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router  /v1/download-jobs/{id} [get]
func (h *DownloadJobs) Get(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.DownloadJobs.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, out)
}

// Cancel download job
// @Summary Cancel download job
// @Tags    download-jobs
// @Produce json
// @Param   id path string true "Job ID (uuid)"
// @Success 200 {object} dbgen.DownloadJob
// @Failure 400 {object} map[string]string
// @Router  /v1/download-jobs/{id} [delete]
func (h *DownloadJobs) Cancel(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.DownloadJobs.Cancel(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// List download jobs for a movie (TMDB ID)
// @Summary List movie download jobs
// @Tags    download-jobs
// @Produce json
// @Param   id path int true "Movie ID (TMDB ID)"
// @Success 200 {array} dbgen.DownloadJob
// @Failure 400 {object} map[string]string
// @Router  /v1/movie/{id}/download-jobs [get]
func (h *DownloadJobs) ListForMovie(c echo.Context) error {
	movieID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid movie ID"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.DownloadJobs.ListByMovie(ctx, movieID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list"})
	}
	return c.JSON(http.StatusOK, out)
}


