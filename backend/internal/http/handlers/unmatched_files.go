package handlers

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type UnmatchedFiles struct{ svc *service.Services }

func NewUnmatchedFiles(s *service.Services) *UnmatchedFiles { return &UnmatchedFiles{svc: s} }

func (h *UnmatchedFiles) RegisterProtected(v1 *echo.Group) {
	v1.GET("/unmatched-files", h.List)
	v1.GET("/unmatched-files/:id", h.Get)
	v1.POST("/unmatched-files/:id/match", h.Match)
	v1.POST("/unmatched-files/:id/dismiss", h.Dismiss)
	v1.POST("/unmatched-files/:id/refresh", h.Refresh)
}

// UnmatchedFileResponse for Swagger
type unmatchedFileSwagger struct {
	ID               string                           `json:"id"`
	LibraryID        string                           `json:"libraryId"`
	Path             string                           `json:"path"`
	FileSize         *int64                           `json:"fileSize,omitempty"`
	DiscoveredAt     string                           `json:"discoveredAt"`
	SuggestedMatches []service.SuggestedMatch `json:"suggestedMatches,omitempty"`
}

// ListResponse for Swagger
type unmatchedFilesListResponse struct {
	Items      []unmatchedFileSwagger `json:"items"`
	TotalCount int64                  `json:"totalCount"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"pageSize"`
}

// MatchRequest payload
type UnmatchedFileMatchRequest struct {
	TmdbID  int64  `json:"tmdbId"`
	Type    string `json:"type"`
	Season  *int   `json:"season,omitempty"`
	Episode *int   `json:"episode,omitempty"`
}

// List unmatched files
// @Summary List unmatched files
// @Description List files that couldn't be auto-matched to media items
// @Tags    unmatched-files
// @Produce json
// @Param   libraryId query string false "Filter by library ID"
// @Param   page query int false "Page number (default 1)"
// @Param   pageSize query int false "Page size (default 20)"
// @Success 200 {object} handlers.unmatchedFilesListResponse
// @Router  /v1/unmatched-files [get]
func (h *UnmatchedFiles) List(c echo.Context) error {
	ctx := c.Request().Context()

	params := service.ListParams{
		Page:     1,
		PageSize: 20,
	}

	if pageStr := c.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if pageSizeStr := c.QueryParam("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			params.PageSize = pageSize
		}
	}

	if libIDStr := c.QueryParam("libraryId"); libIDStr != "" {
		var libID pgtype.UUID
		if err := libID.Scan(libIDStr); err == nil {
			params.LibraryID = &libID
		}
	}

	result, err := h.svc.UnmatchedFiles.List(ctx, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list"})
	}

	return c.JSON(http.StatusOK, result)
}

// Get unmatched file
// @Summary Get unmatched file details
// @Tags    unmatched-files
// @Produce json
// @Param   id path string true "Unmatched file ID"
// @Success 200 {object} handlers.unmatchedFileSwagger
// @Failure 404 {object} map[string]string
// @Router  /v1/unmatched-files/{id} [get]
func (h *UnmatchedFiles) Get(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	file, err := h.svc.UnmatchedFiles.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, file)
}

// Match unmatched file to media item
// @Summary Match unmatched file to media item
// @Description Manually match an unmatched file to a specific media item
// @Tags    unmatched-files
// @Accept  json
// @Produce json
// @Param   id path string true "Unmatched file ID"
// @Param   payload body handlers.UnmatchedFileMatchRequest true "Match request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router  /v1/unmatched-files/{id}/match [post]
func (h *UnmatchedFiles) Match(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req UnmatchedFileMatchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	if req.TmdbID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tmdbId is required"})
	}

	if req.Type != "movie" && req.Type != "series" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "type must be 'movie' or 'series'"})
	}

	ctx := c.Request().Context()
	mediaFile, err := h.svc.UnmatchedFiles.Match(ctx, id, service.MatchRequest{
		TmdbID:  req.TmdbID,
		Type:    req.Type,
		Season:  req.Season,
		Episode: req.Episode,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"mediaFileId": mediaFile.ID.String(),
		"path":        mediaFile.Path,
	})
}

// Dismiss unmatched file
// @Summary Dismiss unmatched file
// @Description Mark an unmatched file as dismissed (resolved without matching)
// @Tags    unmatched-files
// @Param   id path string true "Unmatched file ID"
// @Success 204 {string} string ""
// @Failure 404 {object} map[string]string
// @Router  /v1/unmatched-files/{id}/dismiss [post]
func (h *UnmatchedFiles) Dismiss(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	if err := h.svc.UnmatchedFiles.Dismiss(ctx, id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Refresh match suggestions
// @Summary Refresh match suggestions
// @Description Regenerate match suggestions for an unmatched file
// @Tags    unmatched-files
// @Produce json
// @Param   id path string true "Unmatched file ID"
// @Success 200 {object} handlers.unmatchedFileSwagger
// @Failure 404 {object} map[string]string
// @Router  /v1/unmatched-files/{id}/refresh [post]
func (h *UnmatchedFiles) Refresh(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	file, err := h.svc.UnmatchedFiles.RefreshSuggestions(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, file)
}
