package versioninfo

import "os"

type BuildInfo struct {
	Version    string            `json:"version"`
	Commit     string            `json:"commit,omitempty"`
	BuildDate  string            `json:"buildDate,omitempty"`
	Components map[string]string `json:"components,omitempty"`
}

// Get returns build information from environment variables
func Get() BuildInfo {
	version := os.Getenv("ARRFLIX_VERSION")
	if version == "" {
		version = "dev"
	}

	info := BuildInfo{
		Version:    version,
		Commit:     os.Getenv("ARRFLIX_COMMIT"),
		BuildDate:  os.Getenv("ARRFLIX_BUILD_DATE"),
		Components: make(map[string]string),
	}

	// Add Prowlarr version if available
	if prowlarr := os.Getenv("PROWLARR_VERSION"); prowlarr != "" {
		info.Components["prowlarr"] = prowlarr
	}

	return info
}

// IsDev returns true if this is a dev build
func (b BuildInfo) IsDev() bool {
	return b.Version == "dev"
}

// IsEdge returns true if this is an edge build
func (b BuildInfo) IsEdge() bool {
	return b.Version == "edge"
}

// IsPrerelease returns true if this is a prerelease build (contains -)
func (b BuildInfo) IsPrerelease() bool {
	return !b.IsDev() && !b.IsEdge() && containsDash(b.Version)
}

// IsStable returns true if this is a stable release
func (b BuildInfo) IsStable() bool {
	return !b.IsDev() && !b.IsEdge() && !b.IsPrerelease()
}

func containsDash(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '-' {
			return true
		}
	}
	return false
}
