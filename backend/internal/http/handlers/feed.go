package handlers

import (
	"net/http"

	_ "github.com/kyleaupton/arrflix/internal/model"
	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Feed struct{ svc *service.Services }

func NewFeed(s *service.Services) *Feed { return &Feed{svc: s} }

func (h *Feed) RegisterProtected(v1 *echo.Group) {
	v1.GET("/home", h.GetFeed)
}

// GetFeed
// @Summary Get home feed
// @Tags    feed
// @Produce json
// @Success 200 {object} model.Feed
// @Router  /v1/home [get]
func (h *Feed) GetFeed(c echo.Context) error {
	ctx := c.Request().Context()
	feed, err := h.svc.Feed.GetFeed(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get feed"})
	}
	return c.JSON(http.StatusOK, feed)
}
