package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/Arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type NameTemplates struct{ svc *service.Services }

func NewNameTemplates(s *service.Services) *NameTemplates { return &NameTemplates{svc: s} }

func (h *NameTemplates) RegisterProtected(v1 *echo.Group) {
	v1.GET("/name-templates", h.List)
	v1.POST("/name-templates", h.Create)
	v1.GET("/name-templates/:id", h.Get)
	v1.PUT("/name-templates/:id", h.Update)
	v1.DELETE("/name-templates/:id", h.Delete)
	v1.GET("/name-templates/default/:type", h.GetDefault)
}

// nameTemplateSwagger is a minimal type used only for Swagger schemas.
// It mirrors fields of dbgen.NameTemplate without importing it here.
type nameTemplateSwagger struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Type                 string  `json:"type"`
	Template             string  `json:"template"`
	SeriesShowTemplate   *string `json:"series_show_template"`
	SeriesSeasonTemplate *string `json:"series_season_template"`
	MovieDirTemplate     *string `json:"movie_dir_template"`
	Default              bool    `json:"default"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

// NameTemplateCreateRequest payload
type NameTemplateCreateRequest struct {
	Name                 string  `json:"name"`
	Type                 string  `json:"type"`
	Template             string  `json:"template"`
	SeriesShowTemplate   *string `json:"series_show_template"`
	SeriesSeasonTemplate *string `json:"series_season_template"`
	MovieDirTemplate     *string `json:"movie_dir_template"`
	Default              bool    `json:"default"`
}

// NameTemplateUpdateRequest payload
type NameTemplateUpdateRequest struct {
	Name                 string  `json:"name"`
	Type                 string  `json:"type"`
	Template             string  `json:"template"`
	SeriesShowTemplate   *string `json:"series_show_template"`
	SeriesSeasonTemplate *string `json:"series_season_template"`
	MovieDirTemplate     *string `json:"movie_dir_template"`
	Default              bool    `json:"default"`
}

// List name templates
// @Summary List name templates
// @Tags    name-templates
// @Produce json
// @Success 200 {array} handlers.nameTemplateSwagger
// @Router  /v1/name-templates [get]
func (h *NameTemplates) List(c echo.Context) error {
	ctx := c.Request().Context()
	out, err := h.svc.NameTemplates.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list"})
	}
	return c.JSON(http.StatusOK, out)
}

// Create name template
// @Summary Create name template
// @Tags    name-templates
// @Accept  json
// @Produce json
// @Param   payload body handlers.NameTemplateCreateRequest true "Create name template"
// @Success 201 {object} handlers.nameTemplateSwagger
// @Failure 400 {object} map[string]string
// @Router  /v1/name-templates [post]
func (h *NameTemplates) Create(c echo.Context) error {
	var req NameTemplateCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	ctx := c.Request().Context()
	template, err := h.svc.NameTemplates.Create(ctx, req.Name, req.Type, req.Template, req.SeriesShowTemplate, req.SeriesSeasonTemplate, req.MovieDirTemplate, req.Default)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, template)
}

// Get name template
// @Summary Get name template
// @Tags    name-templates
// @Produce json
// @Success 200 {object} handlers.nameTemplateSwagger
// @Router  /v1/name-templates/{id} [get]
func (h *NameTemplates) Get(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	template, err := h.svc.NameTemplates.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, template)
}

// Update name template
// @Summary Update name template
// @Tags    name-templates
// @Accept  json
// @Produce json
// @Param   id path string true "Name Template ID"
// @Param   payload body handlers.NameTemplateUpdateRequest true "Update name template"
// @Success 200 {object} handlers.nameTemplateSwagger
// @Failure 400 {object} map[string]string
// @Router  /v1/name-templates/{id} [put]
func (h *NameTemplates) Update(c echo.Context) error {
	var req NameTemplateUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	template, err := h.svc.NameTemplates.Update(ctx, id, req.Name, req.Type, req.Template, req.SeriesShowTemplate, req.SeriesSeasonTemplate, req.MovieDirTemplate, req.Default)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, template)
}

// Delete name template
// @Summary Delete name template
// @Tags    name-templates
// @Param   id path string true "Name Template ID"
// @Success 204 {string} string ""
// @Router  /v1/name-templates/{id} [delete]
func (h *NameTemplates) Delete(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	if err := h.svc.NameTemplates.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete"})
	}
	return c.NoContent(http.StatusNoContent)
}

// Get default name template
// @Summary Get default name template by type
// @Tags    name-templates
// @Produce json
// @Param   type path string true "Template type (movie or series)"
// @Success 200 {object} handlers.nameTemplateSwagger
// @Failure 404 {object} map[string]string
// @Router  /v1/name-templates/default/{type} [get]
func (h *NameTemplates) GetDefault(c echo.Context) error {
	typ := c.Param("type")
	if typ != "movie" && typ != "series" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "type must be 'movie' or 'series'"})
	}
	ctx := c.Request().Context()
	template, err := h.svc.NameTemplates.GetDefault(ctx, typ)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, template)
}
