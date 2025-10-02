package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyleaupton/snaggle/backend/internal/config"
	"github.com/kyleaupton/snaggle/backend/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Auth struct {
    cfg  config.Config
    log  zerolog.Logger
    pool *pgxpool.Pool
    svc  *service.Services
}

func NewAuth(cfg config.Config, log zerolog.Logger, pool *pgxpool.Pool, svc *service.Services) *Auth {
    return &Auth{cfg: cfg, log: log, pool: pool, svc: svc}
}

type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type loginResponse struct {
    Token string `json:"token"`
}

func (h *Auth) Login(c echo.Context) error {
    var req loginRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
    }
    if req.Email == "" || req.Password == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "email and password required"})
    }

    ctx := c.Request().Context()
    signed, err := h.svc.Auth.Login(ctx, req.Email, req.Password)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
    }
    return c.JSON(http.StatusOK, loginResponse{Token: signed})
}

type meResponse struct {
    ID          string  `json:"id"`
    Email       *string `json:"email"`
    DisplayName *string `json:"display_name"`
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


