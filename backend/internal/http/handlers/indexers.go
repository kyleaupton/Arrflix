package handlers

import (
	"net/http"
	"strconv"

	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
	"golift.io/starr/prowlarr"
)

type Indexers struct{ svc *service.Services }

func NewIndexers(s *service.Services) *Indexers { return &Indexers{svc: s} }

func (h *Indexers) RegisterProtected(v1 *echo.Group) {
	v1.GET("/indexers/configured", h.ListConfigured)
	v1.GET("/indexers/schema", h.GetSchema)
	v1.GET("/indexer/:id", h.GetIndexer)
	v1.POST("/indexer", h.SaveConfig)
	v1.DELETE("/indexer/:id", h.Delete)
}

// ListConfigured returns only configured indexers
// @Summary List configured indexers
// @Tags    indexers
// @Produce json
// @Success 200 {array} model.IndexerOutput
// @Router  /v1/indexers/configured [get]
func (h *Indexers) ListConfigured(c echo.Context) error {
	ctx := c.Request().Context()
	indexers, err := h.svc.Indexer.ListConfiguredIndexers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list configured indexers"})
	}
	return c.JSON(http.StatusOK, indexers)
}

// ListUnconfigured returns the schema of the indexers
// @Summary Get schema of indexers
// @Tags    indexers
// @Produce json
// @Success 200 {array} model.IndexerSchema
// @Router  /v1/indexers/schema [get]
func (h *Indexers) GetSchema(c echo.Context) error {
	ctx := c.Request().Context()
	schema, err := h.svc.Indexer.GetSchema(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list unconfigured indexers"})
	}
	return c.JSON(http.StatusOK, schema)
}

// Get returns a specific indexer by ID
// @Summary Get indexer by ID
// @Tags    indexers
// @Produce json
// @Param   id path string true "Indexer ID"
// @Success 200 {object} model.IndexerOutput
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/indexers/{id} [get]
func (h *Indexers) GetIndexer(c echo.Context) error {
	indexerID := c.Param("id")
	if indexerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexer ID required"})
	}

	indexerIDInt, err := strconv.ParseInt(indexerID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid indexer ID"})
	}

	ctx := c.Request().Context()
	indexer, err := h.svc.Indexer.GetIndexer(ctx, indexerIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get indexer"})
	}

	return c.JSON(http.StatusOK, indexer)
}

// SaveConfig saves the configuration for a specific indexer
// @Summary Save indexer configuration
// @Tags    indexers
// @Accept  json
// @Produce json
// @Param   payload body model.IndexerInput true "Save indexer"
// @Success 200 {object} model.IndexerOutput
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/indexer [post]
func (h *Indexers) SaveConfig(c echo.Context) error {
	ctx := c.Request().Context()
	var config *prowlarr.IndexerInput
	if err := c.Bind(&config); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	res, err := h.svc.Indexer.SaveIndexerConfig(ctx, config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save indexer config"})
	}

	return c.JSON(http.StatusOK, res)
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
