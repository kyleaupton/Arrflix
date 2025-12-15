package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type DownloadCandidates struct{ svc *service.Services }

func NewDownloadCandidates(s *service.Services) *DownloadCandidates {
	return &DownloadCandidates{svc: s}
}

func (h *DownloadCandidates) RegisterProtected(v1 *echo.Group) {
	v1.GET("/movie/:id/download-candidates", h.GetDownloadCandidates)
	v1.POST("/movie/:id/enqueue-candidate", h.EnqueueCandidate)
}

// GetDownloadCandidates searches for download candidates for a movie
// @Summary Get download candidates for a movie
// @Tags    download-candidates
// @Produce json
// @Param   id path int true "Movie ID (TMDB ID)"
// @Success 200 {array} model.DownloadCandidate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/movie/{id}/download-candidates [get]
func (h *DownloadCandidates) GetDownloadCandidates(c echo.Context) error {
	movieID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid movie ID"})
	}

	ctx := c.Request().Context()
	candidates, err := h.svc.DownloadCandidates.SearchDownloadCandidates(ctx, movieID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, candidates)
}

// EnqueueCandidateRequest is the request body for enqueueing a candidate
type EnqueueCandidateRequest struct {
	IndexerID int64  `json:"indexerId"`
	GUID      string `json:"guid"`
}

// EnqueueCandidate enqueues a download candidate through the policy engine
// @Summary Enqueue a download candidate
// @Tags    download-candidates
// @Accept  json
// @Produce json
// @Param   id path int true "Movie ID (TMDB ID)"
// @Param   payload body handlers.EnqueueCandidateRequest true "Enqueue request"
// @Success 200 {object} model.Plan
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/movie/{id}/enqueue-candidate [post]
func (h *DownloadCandidates) EnqueueCandidate(c echo.Context) error {
	movieID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid movie ID"})
	}

	var req EnqueueCandidateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	if req.IndexerID == 0 || req.GUID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexerId and guid are required"})
	}

	ctx := c.Request().Context()
	plan, err := h.svc.DownloadCandidates.EnqueueCandidate(ctx, movieID, req.IndexerID, req.GUID)
	if err != nil {
		// Check if it's a not found error (cache expired)
		if errors.Is(err, service.ErrCandidateNotFound) || errors.Is(err, service.ErrCandidateExpired) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, plan)
}
