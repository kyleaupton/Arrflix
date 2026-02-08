package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/arrflix/internal/config"
	"github.com/kyleaupton/arrflix/internal/logger"
	"github.com/kyleaupton/arrflix/internal/plex"
	"github.com/kyleaupton/arrflix/internal/service"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	cfg  config.Config
	log  *logger.Logger
	pool *pgxpool.Pool
	svc  *service.Services
}

func (h *Auth) RegisterPublic(v1 *echo.Group) {
	v1.POST("/auth/login", h.Login)
	v1.POST("/auth/signup", h.Signup)
	v1.GET("/auth/plex/start", h.PlexStart)
	v1.POST("/auth/plex/exchange", h.PlexExchange)
}

func (h *Auth) RegisterProtected(v1 *echo.Group) {
	v1.GET("/auth/me", h.Me)
}

func NewAuth(cfg config.Config, log *logger.Logger, pool *pgxpool.Pool, svc *service.Services) *Auth {
	return &Auth{cfg: cfg, log: log, pool: pool, svc: svc}
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token" validate:"required"`
}

// @Summary Login
// @Tags    auth
// @Accept  json
// @Produce json
// @Param   payload body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Router  /v1/auth/login [post]
func (h *Auth) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if req.Login == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "login and password required"})
	}

	ctx := c.Request().Context()
	signed, err := h.svc.Auth.Login(ctx, req.Login, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}
	return c.JSON(http.StatusOK, LoginResponse{Token: signed})
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignupResponse struct {
	Success bool `json:"success"`
}

// @Summary  Signup
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    payload body SignupRequest true "Signup request"
// @Success  201 {object} SignupResponse
// @Failure  400 {object} map[string]string
// @Failure  403 {object} map[string]string
// @Failure  409 {object} map[string]string
// @Router   /v1/auth/signup [post]
func (h *Auth) Signup(c echo.Context) error {
	var req SignupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "email, username, and password are required"})
	}

	ctx := c.Request().Context()

	// Check signup strategy
	all, err := h.svc.Settings.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read settings"})
	}
	strategy, _ := all["auth.signup_strategy"].(string)
	if strategy == "" {
		strategy = "invite_only"
	}

	if strategy == "invite_only" {
		if err := h.svc.Invites.CheckAndClaim(ctx, req.Email); err != nil {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
	}

	_, err = h.svc.Users.Create(ctx, req.Email, req.Username, req.Password, "user", true)
	if err != nil {
		if err.Error() == "email or username already exists" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, SignupResponse{Success: true})
}

// PlexStart initiates the Plex SSO flow by creating a PIN and redirecting to Plex.
func (h *Auth) PlexStart(c echo.Context) error {
	redirectURI := c.QueryParam("redirect_uri")
	if redirectURI == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "redirect_uri required"})
	}

	pc := plex.NewClient()
	pin, err := pc.CreatePin()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create plex pin"})
	}

	// Build forward URL: append pinId to the frontend callback
	sep := "?"
	if strings.Contains(redirectURI, "?") {
		sep = "&"
	}
	forwardURL := fmt.Sprintf("%s%spinId=%d", redirectURI, sep, pin.ID)

	authURL := plex.AuthURL(pin.Code, forwardURL)
	return c.Redirect(http.StatusFound, authURL)
}

type PlexExchangeRequest struct {
	PinID int `json:"pin_id" validate:"required"`
}

type PlexExchangeResponse struct {
	Token string `json:"token" validate:"required"`
}

// @Summary  Exchange Plex PIN for JWT
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    payload body PlexExchangeRequest true "Plex exchange request"
// @Success  200 {object} PlexExchangeResponse
// @Failure  400 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Failure  403 {object} map[string]string
// @Failure  409 {object} map[string]string
// @Router   /v1/auth/plex/exchange [post]
func (h *Auth) PlexExchange(c echo.Context) error {
	var req PlexExchangeRequest
	if err := c.Bind(&req); err != nil || req.PinID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "pin_id required"})
	}

	pc := plex.NewClient()

	// Check if the PIN has been claimed
	pin, err := pc.CheckPin(req.PinID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to check plex pin"})
	}
	if pin.AuthToken == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "plex authorization not completed"})
	}

	// Get Plex user info
	plexUser, err := pc.GetUser(pin.AuthToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get plex user info"})
	}

	// Marshal raw Plex user data for storage
	raw, _ := json.Marshal(plexUser)

	plexSubject := strconv.Itoa(plexUser.ID)
	ctx := c.Request().Context()

	signed, err := h.svc.Auth.LoginWithPlex(ctx, plexSubject, plexUser.Email, plexUser.Username, pin.AuthToken, raw)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "no invite") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": msg})
		}
		if strings.Contains(msg, "already exists") {
			return c.JSON(http.StatusConflict, map[string]string{"error": msg})
		}
		if strings.Contains(msg, "disabled") {
			return c.JSON(http.StatusForbidden, map[string]string{"error": msg})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "login failed"})
	}

	return c.JSON(http.StatusOK, PlexExchangeResponse{Token: signed})
}

type MeResponse struct {
	ID       string  `json:"id" validate:"required"`
	Email    *string `json:"email" validate:"required"`
	Username *string `json:"username" validate:"required"`
}

func (h *Auth) Me(c echo.Context) error {
	authz := c.Request().Header.Get("Authorization")
	if authz == "" || len(authz) < 8 || authz[:7] != "Bearer " {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
	}
	raw := authz[7:]
	tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "bad token method")
		}
		return []byte(h.cfg.JWTSecret), nil
	})
	if err != nil || !tok.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}
	sub, _ := claims["sub"].(string)
	if sub == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token subject"})
	}

	// For now, just echo back minimal info from token; optionally you can fetch from DB later
	return c.JSON(http.StatusOK, map[string]interface{}{
		"sub":   sub,
		"email": claims["email"],
		"name":  claims["name"],
	})
}
