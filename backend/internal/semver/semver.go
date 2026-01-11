package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a semantic version
type Version struct {
	Major int
	Minor int
	Patch int
}

// Parse parses a semver string (e.g., "1.4.2" or "v1.4.2")
func Parse(s string) (*Version, error) {
	// Strip leading 'v'
	s = strings.TrimPrefix(s, "v")

	// Split by '-' to remove prerelease suffix
	parts := strings.Split(s, "-")
	s = parts[0]

	// Split by '.'
	components := strings.Split(s, ".")
	if len(components) != 3 {
		return nil, fmt.Errorf("invalid semver format: %s", s)
	}

	major, err := strconv.Atoi(components[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %w", err)
	}

	minor, err := strconv.Atoi(components[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %w", err)
	}

	patch, err := strconv.Atoi(components[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %w", err)
	}

	return &Version{Major: major, Minor: minor, Patch: patch}, nil
}

// Compare returns:
//
//	-1 if v < other
//	 0 if v == other
//	 1 if v > other
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// LessThan returns true if v < other
func (v *Version) LessThan(other *Version) bool {
	return v.Compare(other) < 0
}
