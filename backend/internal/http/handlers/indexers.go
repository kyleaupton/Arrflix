package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	v1.PUT("/indexer/:id/toggle", h.Toggle)
	v1.POST("/indexer/action/:name", h.Action)
	v1.POST("/indexer/:id/test", h.TestSaved)
	v1.POST("/indexer/test", h.TestUnsaved)
	v1.POST("/indexers/testall", h.TestAll)
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
// @Router  /v1/indexer/{id} [delete]
func (h *Indexers) Delete(c echo.Context) error {
	fmt.Println("Delete indexer")
	indexerID := c.Param("id")
	if indexerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexer ID required"})
	}

	indexerIDInt, err := strconv.ParseInt(indexerID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid indexer ID"})
	}

	ctx := c.Request().Context()
	if err := h.svc.Indexer.DeleteIndexer(ctx, indexerIDInt); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete indexer"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Toggle toggles the enable state of an indexer
// @Summary Toggle indexer enable state
// @Tags    indexers
// @Param   id path string true "Indexer ID"
// @Success 200 {object} model.IndexerOutput
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/indexer/{id}/toggle [put]
func (h *Indexers) Toggle(c echo.Context) error {
	indexerID := c.Param("id")
	if indexerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "indexer ID required"})
	}

	indexerIDInt, err := strconv.ParseInt(indexerID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid indexer ID"})
	}

	ctx := c.Request().Context()
	result, err := h.svc.Indexer.ToggleIndexer(ctx, indexerIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to toggle indexer"})
	}

	return c.JSON(http.StatusOK, result)
}

// Action performs an action on an indexer
// @Summary Perform an action on an indexer
// @Tags    indexers
// @Param   name path string true "Action name"
// @Param   payload body model.IndexerDefinition true "Action input"
// @Success 200 {object} any
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/indexer/action/{name} [post]
func (h *Indexers) Action(c echo.Context) error {
	actionName := c.Param("name")
	if actionName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "action name required"})
	}

	var input interface{}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	ctx := c.Request().Context()
	action, err := h.svc.Indexer.Action(ctx, actionName, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to perform action"})
	}

	return c.JSON(http.StatusOK, action)
}

// TestSaved tests a saved indexer configuration
// @Summary Test saved indexer
// @Tags    indexers
// @Param   id path int64 true "Indexer ID"
// @Success 200 {object} model.IndexerTestResult
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/indexer/{id}/test [post]
func (h *Indexers) TestSaved(c echo.Context) error {
	idStr := c.Param("id")
	indexerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid indexer ID"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 15*time.Second)
	defer cancel()

	result, err := h.svc.Indexer.TestIndexerByID(ctx, indexerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// TestUnsaved tests an unsaved indexer configuration
// @Summary Test unsaved indexer configuration
// @Tags    indexers
// @Accept  json
// @Param   config body model.IndexerInput true "Indexer config"
// @Success 200 {object} model.IndexerTestResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router  /v1/indexer/test [post]
func (h *Indexers) TestUnsaved(c echo.Context) error {
	var config prowlarr.IndexerInput
	if err := c.Bind(&config); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 15*time.Second)
	defer cancel()

	result, err := h.svc.Indexer.TestIndexer(ctx, &config)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// TestAll tests all configured indexers
// @Summary Test all indexers
// @Tags    indexers
// @Success 200 {array} model.IndexerBatchTestResult
// @Failure 500 {object} map[string]string
// @Router  /v1/indexers/testall [post]
func (h *Indexers) TestAll(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 45*time.Second)
	defer cancel()

	results, err := h.svc.Indexer.TestAllIndexers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}
