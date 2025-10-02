package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kyleaupton/snaggle/backend/internal/repo"
)

type AuthService struct {
    repo      *repo.Repository
    jwtSecret string
}

func NewAuthService(r *repo.Repository, cfg *config) *AuthService {
    return &AuthService{repo: r, jwtSecret: cfg.jwtSecret}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
    user, err := s.repo.GetUserByEmail(ctx, email)
    if err != nil { return "", err }

    sum := sha256.Sum256([]byte(password))
    if user.PasswordHash == nil || *user.PasswordHash != hex.EncodeToString(sum[:]) {
        return "", ErrInvalidCredentials
    }

    claims := jwt.MapClaims{
        "sub": user.ID.String(),
        "email": user.Email,
        "name": user.DisplayName,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
        "iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

var ErrInvalidCredentials = jwt.ErrInvalidKey