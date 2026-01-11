package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kyleaupton/Arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Users struct{ svc *service.Services }

func NewUsers(s *service.Services) *Users { return &Users{svc: s} }

func (h *Users) RegisterProtected(v1 *echo.Group) {
	// Admin endpoints
	v1.GET("/users", h.List)
	v1.POST("/users", h.Create)
	v1.GET("/users/:id", h.Get)
	v1.PUT("/users/:id", h.Update)
	v1.DELETE("/users/:id", h.Delete)
	v1.PUT("/users/:id/password", h.UpdatePassword)
	v1.PUT("/users/:id/role", h.AssignRole)

	// Profile endpoints (current user)
	v1.GET("/auth/profile", h.GetProfile)
	v1.PUT("/auth/profile", h.UpdateProfile)
	v1.PUT("/auth/profile/password", h.UpdateProfilePassword)

	// Utility
	v1.GET("/roles", h.ListRoles)
}

// Swagger models
type userSwagger struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	DisplayName string   `json:"display_name"`
	IsActive    bool     `json:"is_active"`
	Roles       []string `json:"roles"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type UserCreateRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	IsActive    bool   `json:"is_active"`
}

type UserUpdateRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	IsActive    bool   `json:"is_active"`
}

type UserPasswordUpdateRequest struct {
	Password string `json:"password"`
}

type UserRoleUpdateRequest struct {
	Role string `json:"role"`
}

// List users
// @Summary List users
// @Tags    users
// @Produce json
// @Success 200 {array} handlers.userSwagger
// @Router  /v1/users [get]
func (h *Users) List(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.svc.Users.List(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
	}
	return c.JSON(http.StatusOK, users)
}

// Create user
// @Summary Create user
// @Tags    users
// @Accept  json
// @Produce json
// @Param   payload body handlers.UserCreateRequest true "Create user"
// @Success 201 {object} handlers.userSwagger
// @Failure 400 {object} map[string]string
// @Router  /v1/users [post]
func (h *Users) Create(c echo.Context) error {
	var req UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	ctx := c.Request().Context()
	user, err := h.svc.Users.Create(ctx, req.Email, req.DisplayName, req.Password, req.Role, req.IsActive)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// Get user
// @Summary Get user
// @Tags    users
// @Produce json
// @Param   id path string true "User ID"
// @Success 200 {object} handlers.userSwagger
// @Router  /v1/users/{id} [get]
func (h *Users) Get(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	user, err := h.svc.Users.Get(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// Update user
// @Summary Update user
// @Tags    users
// @Accept  json
// @Produce json
// @Param   id path string true "User ID"
// @Param   payload body handlers.UserUpdateRequest true "Update user"
// @Success 200 {object} handlers.userSwagger
// @Failure 400 {object} map[string]string
// @Router  /v1/users/{id} [put]
func (h *Users) Update(c echo.Context) error {
	var req UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	user, err := h.svc.Users.Update(ctx, id, req.Email, req.DisplayName, req.IsActive)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// Delete user
// @Summary Delete user
// @Tags    users
// @Param   id path string true "User ID"
// @Success 204 {string} string ""
// @Router  /v1/users/{id} [delete]
func (h *Users) Delete(c echo.Context) error {
	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	if err := h.svc.Users.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdatePassword changes a user's password
// @Summary Update user password
// @Tags    users
// @Accept  json
// @Param   id path string true "User ID"
// @Param   payload body handlers.UserPasswordUpdateRequest true "Update password"
// @Success 204 {string} string ""
// @Router  /v1/users/{id}/password [put]
func (h *Users) UpdatePassword(c echo.Context) error {
	var req UserPasswordUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	if err := h.svc.Users.UpdatePassword(ctx, id, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// AssignRole assigns a role to a user
// @Summary Assign role to user
// @Tags    users
// @Accept  json
// @Param   id path string true "User ID"
// @Param   payload body handlers.UserRoleUpdateRequest true "Assign role"
// @Success 204 {string} string ""
// @Router  /v1/users/{id}/role [put]
func (h *Users) AssignRole(c echo.Context) error {
	var req UserRoleUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	var id pgtype.UUID
	if err := id.Scan(c.Param("id")); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	ctx := c.Request().Context()
	if err := h.svc.Users.AssignRole(ctx, id, req.Role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetProfile returns current user's profile
// @Summary Get current user profile
// @Tags    auth
// @Produce json
// @Success 200 {object} handlers.userSwagger
// @Router  /v1/auth/profile [get]
func (h *Users) GetProfile(c echo.Context) error {
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

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	ctx := c.Request().Context()
	user, err := h.svc.Users.Get(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateProfile updates current user's profile
// @Summary Update current user profile
// @Tags    auth
// @Accept  json
// @Produce json
// @Param   payload body handlers.UserUpdateRequest true "Update profile"
// @Success 200 {object} handlers.userSwagger
// @Router  /v1/auth/profile [put]
func (h *Users) UpdateProfile(c echo.Context) error {
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

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	var req UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	// Users cannot change their own active status via profile
	ctx := c.Request().Context()
	currentUser, err := h.svc.Users.Get(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	user, err := h.svc.Users.Update(ctx, userID, req.Email, req.DisplayName, currentUser.IsActive)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateProfilePassword changes current user's password
// @Summary Update current user password
// @Tags    auth
// @Accept  json
// @Param   payload body handlers.UserPasswordUpdateRequest true "Update password"
// @Success 204 {string} string ""
// @Router  /v1/auth/profile/password [put]
func (h *Users) UpdateProfilePassword(c echo.Context) error {
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

	var userID pgtype.UUID
	if err := userID.Scan(userIDStr); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	var req UserPasswordUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	ctx := c.Request().Context()
	if err := h.svc.Users.UpdatePassword(ctx, userID, req.Password); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// ListRoles returns available roles
// @Summary List roles
// @Tags    users
// @Produce json
// @Success 200 {array} dbgen.Role
// @Router  /v1/roles [get]
func (h *Users) ListRoles(c echo.Context) error {
	ctx := c.Request().Context()
	roles, err := h.svc.Users.ListRoles(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list roles"})
	}
	return c.JSON(http.StatusOK, roles)
}
