package handlers

import (
	"net/http"

	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Media struct{ svc *service.Services }

func NewMedia(s *service.Services) *Media { return &Media{svc: s} }

func (h *Media) RegisterProtected(v1 *echo.Group) {
	v1.GET("/library", h.List)
}

// mediaItemSwagger mirrors fields of dbgen.MediaItem for Swagger without importing it here.
type mediaItemSwagger struct {
	ID        string `json:"id"`
	LibraryID string `json:"library_id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Year      *int32 `json:"year"`
	TmdbID    *int32 `json:"tmdb_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// List media items
// @Summary List media items
// @Tags    media
// @Produce json
// @Success 200 {array} handlers.mediaItemSwagger
// @Router  /v1/library [get]
func (h *Media) List(c echo.Context) error {
	ctx := c.Request().Context()
	items, err := h.svc.Media.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list"})
	}
	return c.JSON(http.StatusOK, items)
}
