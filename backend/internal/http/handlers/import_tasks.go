package handlers

import (
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type ImportTasks struct{ svc *service.Services }

func NewImportTasks(s *service.Services) *ImportTasks { return &ImportTasks{svc: s} }

func (h *ImportTasks) RegisterProtected(v1 *echo.Group) {
	v1.GET("/import-tasks", h.List)
	v1.GET("/import-tasks/counts", h.Counts)
	v1.GET("/import-tasks/:id", h.Get)
	v1.GET("/import-tasks/:id/timeline", h.GetTimeline)
	v1.GET("/import-tasks/:id/history", h.GetHistory)
	v1.POST("/import-tasks/:id/reimport", h.Reimport)
	v1.POST("/import-tasks/:id/cancel", h.Cancel)
}

// List import tasks
// @Summary List import tasks
// @Tags    import-tasks
// @Produce json
// @Param   status query string false "Filter by status"
// @Param   limit query int false "Page size" default(50)
// @Param   offset query int false "Page offset" default(0)
// @Success 200 {array} dbgen.ImportTask
// @Router  /v1/import-tasks [get]
func (h *ImportTasks) List(c echo.Context) error {
	ctx := c.Request().Context()
	status := c.QueryParam("status")
	limit := parseIntQuery(c, "limit", 50)
	offset := parseIntQuery(c, "offset", 0)

	var out []any
	var err error

	if status != "" {
		tasks, e := h.svc.ImportTasks.ListByStatus(ctx, status, int32(limit), int32(offset))
		err = e
		for _, t := range tasks {
			out = append(out, t)
		}
	} else {
		tasks, e := h.svc.ImportTasks.List(ctx, int32(limit), int32(offset))
		err = e
		for _, t := range tasks {
			out = append(out, t)
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list"})
	}
	return c.JSON(http.StatusOK, out)
}

// Counts returns import task counts by status
// @Summary Get import task counts by status
// @Tags    import-tasks
// @Produce json
// @Success 200 {object} dbgen.CountImportTasksByStatusRow
// @Router  /v1/import-tasks/counts [get]
func (h *ImportTasks) Counts(c echo.Context) error {
	ctx := c.Request().Context()
	out, err := h.svc.ImportTasks.CountByStatus(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to count"})
	}
	return c.JSON(http.StatusOK, out)
}

// Get import task with details
// @Summary Get import task details
// @Tags    import-tasks
// @Produce json
// @Param   id path string true "Task ID (uuid)"
// @Success 200 {object} dbgen.GetImportTaskWithDetailsRow
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router  /v1/import-tasks/{id} [get]
func (h *ImportTasks) Get(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.ImportTasks.GetWithDetails(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, out)
}

// GetTimeline returns event log for an import task
// @Summary Get import task timeline
// @Tags    import-tasks
// @Produce json
// @Param   id path string true "Task ID (uuid)"
// @Success 200 {array} dbgen.ImportTaskEvent
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router  /v1/import-tasks/{id}/timeline [get]
func (h *ImportTasks) GetTimeline(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.ImportTasks.GetTimeline(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, out)
}

// GetHistory returns the reimport chain for an import task
// @Summary Get import task reimport history
// @Tags    import-tasks
// @Produce json
// @Param   id path string true "Task ID (uuid)"
// @Success 200 {array} dbgen.GetImportTaskHistoryRow
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router  /v1/import-tasks/{id}/history [get]
func (h *ImportTasks) GetHistory(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.ImportTasks.GetHistory(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, out)
}

// Reimport creates a new import task for a completed or failed task
// @Summary Reimport a file
// @Tags    import-tasks
// @Produce json
// @Param   id path string true "Task ID (uuid)"
// @Success 200 {object} dbgen.ImportTask
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router  /v1/import-tasks/{id}/reimport [post]
func (h *ImportTasks) Reimport(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.ImportTasks.Reimport(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// Cancel cancels a pending import task
// @Summary Cancel import task
// @Tags    import-tasks
// @Produce json
// @Param   id path string true "Task ID (uuid)"
// @Success 200 {object} dbgen.ImportTask
// @Failure 400 {object} map[string]string
// @Router  /v1/import-tasks/{id}/cancel [post]
func (h *ImportTasks) Cancel(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	ctx := c.Request().Context()
	out, err := h.svc.ImportTasks.Cancel(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

func parseIntQuery(c echo.Context, key string, def int) int {
	val := c.QueryParam(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}
