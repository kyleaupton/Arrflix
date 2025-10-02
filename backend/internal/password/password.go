package password

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	VersionV1   = "v1"
	AlgoBcrypt  = "bcrypt"
	CurrentCost = 12
)

// Hash returns a versioned password hash string for storage.
// Format: v1:bcrypt:$2a$<cost>$<salt+hash>
func Hash(raw string) (string, error) {
	encoded, err := bcrypt.GenerateFromPassword([]byte(raw), CurrentCost)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{VersionV1, AlgoBcrypt, string(encoded)}, ":"), nil
}

// Verify checks a plaintext password against a stored hash string.
// Supports:
// - v1:bcrypt:<bcrypt>
func Verify(raw, stored string) (bool, error) {
	if stored == "" {
		return false, nil
	}
	stored = strings.TrimSpace(stored)

	// Versioned format
	if strings.HasPrefix(stored, VersionV1+":") {
		parts := strings.SplitN(stored, ":", 3)
		if len(parts) != 3 {
			return false, errors.New("invalid password format")
		}
		algo, payload := strings.ToLower(strings.TrimSpace(parts[1])), strings.TrimSpace(parts[2])
		switch algo {
		case AlgoBcrypt:
			return compareBcrypt(payload, raw), nil
		default:
			return false, errors.New("unsupported password algorithm")
		}
	}

	// Unknown/unsupported format
	return false, nil
}

func compareBcrypt(encoded string, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encoded), []byte(raw)) == nil
}

// NeedsRehash returns true if the stored hash is not in the current preferred format/cost.
func NeedsRehash(stored string) bool {
	if stored == "" {
		return true
	}
	if strings.HasPrefix(stored, VersionV1+":") {
		parts := strings.SplitN(stored, ":", 3)
		if len(parts) != 3 {
			return true
		}
		if parts[1] != AlgoBcrypt {
			return true
		}
		cost, err := bcrypt.Cost([]byte(parts[2]))
		if err != nil {
			return true
		}
		return cost < CurrentCost
	}
	// legacy formats should be upgraded
	if strings.HasPrefix(stored, "$2a$") || strings.HasPrefix(stored, "$2b$") || strings.HasPrefix(stored, "$2y$") {
		cost, err := bcrypt.Cost([]byte(stored))
		if err != nil {
			return true
		}
		return cost < CurrentCost
	}
	return true
}
