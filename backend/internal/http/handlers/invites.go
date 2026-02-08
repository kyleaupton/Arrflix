package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Invites struct{ svc *service.Services }

func NewInvites(s *service.Services) *Invites { return &Invites{svc: s} }

func (h *Invites) RegisterProtected(v1 *echo.Group) {
	v1.GET("/invites", h.List)
	v1.POST("/invites", h.Create)
	v1.DELETE("/invites/:id", h.Delete)
}

type InviteCreateRequest struct {
	Email string `json:"email"`
}

// inviteSwagger is used only for Swagger documentation.
type inviteSwagger struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	InvitedBy string  `json:"invited_by"`
	CreatedAt string  `json:"created_at"`
	ClaimedAt *string `json:"claimed_at"`
}

// List invites
// @Summary List invites
// @Tags    invites
// @Produce json
// @Success 200 {array} handlers.inviteSwagger
// @Router  /v1/invites [get]
func (h *Invites) List(c echo.Context) error {
	ctx := c.Request().Context()
	invites, err := h.svc.Invites.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list invites"})
	}
	return c.JSON(http.StatusOK, invites)
}

// Create invite
// @Summary Create invite
// @Tags    invites
// @Accept  json
// @Produce json
// @Param   payload body handlers.InviteCreateRequest true "Create invite"
// @Success 201 {object} handlers.inviteSwagger
// @Failure 400 {object} map[string]string
// @Router  /v1/invites [post]
func (h *Invites) Create(c echo.Context) error {
	var req InviteCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	claims := c.Get("claims")
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
	}
	userIDStr, ok := claimsMap["sub"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token subject"})
	}
	var invitedBy pgtype.UUID
	if err := invitedBy.Scan(userIDStr); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	ctx := c.Request().Context()
	invite, err := h.svc.Invites.Create(ctx, req.Email, invitedBy)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, invite)
}

// Delete invite
// @Summary Delete invite
// @Tags    invites
// @Param   id path string true "Invite ID"
// @Success 204 {string} string ""
// @Router  /v1/invites/{id} [delete]
func (h *Invites) Delete(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	if err := h.svc.Invites.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
