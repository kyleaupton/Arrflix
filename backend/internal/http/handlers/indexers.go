package handlers

import (
	"net/http"

	"github.com/kyleaupton/snaggle/backend/internal/jackett"
	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
)

type Indexers struct{ svc *service.Services }

func NewIndexers(s *service.Services) *Indexers { return &Indexers{svc: s} }

func (h *Indexers) RegisterProtected(v1 *echo.Group) {
	v1.GET("/indexers", h.ListAll)
	v1.GET("/indexers/configured", h.ListConfigured)
	v1.GET("/indexers/unconfigured", h.ListUnconfigured)
	v1.GET("/indexers/:id", h.Get)
	v1.GET("/indexers/:id/config", h.GetConfig)
	v1.POST("/indexers", h.Create)
	v1.PUT("/indexers/:id/config", h.UpdateConfig)
	v1.DELETE("/indexers/:id", h.Delete)
}

// IndexerCreateRequest represents the request payload for creating an indexer
type IndexerCreateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Enabled     bool                   `json:"enabled"`
	Fields      map[string]interface{} `json:"fields"`
}

// IndexerUpdateRequest represents the request payload for updating an indexer
type IndexerUpdateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Enabled     bool                   `json:"enabled"`
	Fields      map[string]interface{} `json:"fields"`
}

// ListAll returns all indexers (both configured and unconfigured)
// @Summary List all indexers
// @Tags    indexers
// @Produce json
// @Success 200 {array} map[string]any
// @Router  /v1/indexers [get]
func (h *Indexers) ListAll(c echo.Context) error {
	ctx := c.Request().Context()
	indexers, err := h.svc.Indexer.ListAllIndexers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list indexers"})
	}
	return c.JSON(http.StatusOK, indexers)
}

// ListConfigured returns only configured indexers
// @Summary List configured indexers
// @Tags    indexers
// @Produce json
// @Success 200 {array} map[string]any
// @Router  /v1/indexers/configured [get]
func (h *Indexers) ListConfigured(c echo.Context) error {
	ctx := c.Request().Context()
	indexers, err := h.svc.Indexer.IndexersConfigured(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list configured indexers"})
	}
	return c.JSON(http.StatusOK, indexers)
}

// ListUnconfigured returns only unconfigured indexers
// @Summary List unconfigured indexers
// @Tags    indexers
// @Produce json
// @Success 200 {array} map[string]any
// @Router  /v1/indexers/unconfigured [get]
func (h *Indexers) ListUnconfigured(c echo.Context) error {
	ctx := c.Request().Context()
	indexers, err := h.svc.Indexer.IndexersUnconfigured(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list unconfigured indexers"})
	}
	return c.JSON(http.StatusOK, indexers)
}

// Get returns a specific indexer by ID
// @Summary Get indexer by ID
// @Tags    indexers
// @Produce json
// @Param   id path string true "Indexer ID"
// @Success 200 {object} map[string]any
// @Failure 404 {object} map[string]string
// @Router  /v1/indexers/{id} [get]
func (h *Indexers) Get(c echo.Context) error {
	indexerID := c.Param("id")
	if indexerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexer ID required"})
	}

	ctx := c.Request().Context()
	allIndexers, err := h.svc.Indexer.ListAllIndexers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list indexers"})
	}

	// Find the specific indexer
	for _, indexer := range allIndexers {
		if id, ok := indexer["id"].(string); ok && id == indexerID {
			return c.JSON(http.StatusOK, indexer)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"error": "indexer not found"})
}

// GetConfig returns the configuration for a specific indexer
// @Summary Get indexer configuration
// @Tags    indexers
// @Produce json
// @Param   id path string true "Indexer ID"
// @Success 200 {object} jackett.IndexerConfig
// @Failure 404 {object} map[string]string
// @Router  /v1/indexers/{id}/config [get]
func (h *Indexers) GetConfig(c echo.Context) error {
	indexerID := c.Param("id")
	if indexerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexer ID required"})
	}

	ctx := c.Request().Context()
	config, err := h.svc.Indexer.GetIndexerConfig(ctx, indexerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "indexer config not found"})
	}

	return c.JSON(http.StatusOK, config)
}

// Create creates a new indexer
// @Summary Create indexer
// @Tags    indexers
// @Accept  json
// @Produce json
// @Param   payload body handlers.IndexerCreateRequest true "Create indexer"
// @Success 201 {object} jackett.IndexerConfig
// @Failure 400 {object} map[string]string
// @Router  /v1/indexers [post]
func (h *Indexers) Create(c echo.Context) error {
	var req IndexerCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	ctx := c.Request().Context()
	config := &jackett.IndexerConfigRequest{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Enabled:     req.Enabled,
		Fields:      req.Fields,
	}

	indexer, err := h.svc.Indexer.AddIndexer(ctx, config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, indexer)
}

// UpdateConfig updates the configuration for a specific indexer
// @Summary Update indexer configuration
// @Tags    indexers
// @Accept  json
// @Produce json
// @Param   id path string true "Indexer ID"
// @Param   payload body handlers.IndexerUpdateRequest true "Update indexer"
// @Success 200 {object} jackett.IndexerConfig
// @Failure 400 {object} map[string]string
// @Router  /v1/indexers/{id}/config [put]
func (h *Indexers) UpdateConfig(c echo.Context) error {
	indexerID := c.Param("id")
	if indexerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexer ID required"})
	}

	var req IndexerUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	ctx := c.Request().Context()
	config := &jackett.IndexerConfigRequest{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Enabled:     req.Enabled,
		Fields:      req.Fields,
	}

	indexer, err := h.svc.Indexer.UpdateIndexerConfig(ctx, indexerID, config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, indexer)
}

// Delete removes an indexer by ID
// @Summary Delete indexer
// @Tags    indexers
// @Param   id path string true "Indexer ID"
// @Success 204 {string} string ""
// @Failure 400 {object} map[string]string
// @Router  /v1/indexers/{id} [delete]
func (h *Indexers) Delete(c echo.Context) error {
	indexerID := c.Param("id")
	if indexerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexer ID required"})
	}

	ctx := c.Request().Context()
	if err := h.svc.Indexer.DeleteIndexer(ctx, indexerID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete indexer"})
	}

	return c.NoContent(http.StatusNoContent)
}
