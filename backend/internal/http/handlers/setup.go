package handlers

import (
	"net/http"

	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Setup struct {
	svc *service.Services
}

func NewSetup(s *service.Services) *Setup {
	return &Setup{svc: s}
}

// RegisterPublic registers public setup routes (no auth required)
func (h *Setup) RegisterPublic(v1 *echo.Group) {
	v1.GET("/setup/status", h.GetStatus)
	v1.POST("/setup/initialize", h.Initialize)
}

type SetupStatusResponse struct {
	Initialized bool `json:"initialized"`
}

type SetupInitializeRequest struct {
	Email       string `json:"email" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type SetupInitializeResponse struct {
	Success bool `json:"success"`
}

// GetStatus checks if the system has been initialized
// @Summary Get setup status
// @Tags    setup
// @Produce json
// @Success 200 {object} handlers.SetupStatusResponse
// @Router  /v1/setup/status [get]
func (h *Setup) GetStatus(c echo.Context) error {
	ctx := c.Request().Context()
	initialized, err := h.svc.Setup.IsInitialized(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to check status"})
	}
	return c.JSON(http.StatusOK, SetupStatusResponse{Initialized: initialized})
}

// Initialize performs the one-time setup
// @Summary Initialize system
// @Tags    setup
// @Accept  json
// @Produce json
// @Param   payload body handlers.SetupInitializeRequest true "Setup request"
// @Success 200 {object} handlers.SetupInitializeResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string "Already initialized"
// @Router  /v1/setup/initialize [post]
func (h *Setup) Initialize(c echo.Context) error {
	var req SetupInitializeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	ctx := c.Request().Context()
	err := h.svc.Setup.Initialize(ctx, req.Email, req.DisplayName, req.Password)
	if err != nil {
		if err == service.ErrAlreadyInitialized {
			return c.JSON(http.StatusConflict, map[string]string{"error": "system already initialized"})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, SetupInitializeResponse{Success: true})
}
