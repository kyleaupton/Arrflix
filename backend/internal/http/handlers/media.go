package handlers

import (
	"net/http"
	"strconv"

	_ "github.com/jackc/pgx/v5/pgtype"
	_ "github.com/kyleaupton/snaggle/backend/internal/db/sqlc"
	_ "github.com/kyleaupton/snaggle/backend/internal/model"
	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Media struct{ svc *service.Services }

func NewMedia(s *service.Services) *Media { return &Media{svc: s} }

func (h *Media) RegisterProtected(v1 *echo.Group) {
	v1.GET("/library", h.List)
	v1.GET("/movie/:id", h.GetMovie)
	v1.GET("/series/:id", h.GetSeries)
	v1.GET("/person/:id", h.GetPerson)
}

// List media items
// @Summary List media items
// @Tags    media
// @Produce json
// @Success 200 {array} dbgen.MediaItem
// @Router  /v1/library [get]
func (h *Media) List(c echo.Context) error {
	ctx := c.Request().Context()
	items, err := h.svc.Media.ListLibraryItems(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list"})
	}
	return c.JSON(http.StatusOK, items)
}

// GetMovie
// @Summary Get movie (by TMDB id)
// @Tags    media
// @Produce json
// @Param   id path int true "Movie ID"
// @Success 200 {object} model.MovieDetail
// @Router  /v1/movie/{id} [get]
func (h *Media) GetMovie(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	item, err := h.svc.Media.GetMovieDetail(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get movie"})
	}
	return c.JSON(http.StatusOK, item)
}

// GetSeries
// @Summary Get series (by TMDB id)
// @Tags    media
// @Produce json
// @Param   id path int true "Series ID"
// @Success 200 {object} model.SeriesDetail
// @Router  /v1/series/{id} [get]
func (h *Media) GetSeries(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	item, err := h.svc.Media.GetSeriesDetail(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get series"})
	}
	return c.JSON(http.StatusOK, item)
}

// GetPerson
// @Summary Get person (by TMDB id)
// @Tags    media
// @Produce json
// @Param   id path int true "Person ID"
// @Success 200 {object} model.PersonDetail
// @Router  /v1/person/{id} [get]
func (h *Media) GetPerson(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	item, err := h.svc.Media.GetPersonDetail(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get person"})
	}
	return c.JSON(http.StatusOK, item)
}
