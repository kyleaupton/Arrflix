package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	pw "github.com/kyleaupton/Arrflix/internal/password"
	"github.com/kyleaupton/Arrflix/internal/repo"
)

type AuthService struct {
	repo      *repo.Repository
	jwtSecret string
}

func NewAuthService(r *repo.Repository, cfg *cfg) *AuthService {
	return &AuthService{repo: r, jwtSecret: cfg.jwtSecret}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	ok, err := pw.Verify(password, deref(user.PasswordHash))
	if err != nil || !ok {
		return "", ErrInvalidCredentials
	}

	// Opportunistically upgrade hash format/cost
	if pw.NeedsRehash(deref(user.PasswordHash)) {
		if newHash, err := pw.Hash(password); err == nil {
			_ = s.repo.UpdateUserPassword(ctx, user.ID, newHash)
		}
	}

	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"name":  user.DisplayName,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

var ErrInvalidCredentials = jwt.ErrInvalidKey

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
