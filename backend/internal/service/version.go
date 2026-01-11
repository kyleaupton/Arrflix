package service

import (
	"context"
	"time"

	"github.com/kyleaupton/Arrflix/internal/github"
	"github.com/kyleaupton/Arrflix/internal/logger"
	"github.com/kyleaupton/Arrflix/internal/repo"
	"github.com/kyleaupton/Arrflix/internal/semver"
	"github.com/kyleaupton/Arrflix/internal/versioninfo"
)

const (
	GitHubOwner = "kyleaupton"
	GitHubRepo  = "Arrflix"
	UpdateTTL   = 15 * time.Minute
)

type VersionService struct {
	repo   *repo.Repository
	logger *logger.Logger
	gh     *github.Client
	build  versioninfo.BuildInfo
}

func NewVersionService(r *repo.Repository, l *logger.Logger) *VersionService {
	return &VersionService{
		repo:   r,
		logger: l,
		gh:     github.NewClient(GitHubOwner, GitHubRepo),
		build:  versioninfo.Get(),
	}
}

// GetBuildInfo returns current build information
func (s *VersionService) GetBuildInfo() versioninfo.BuildInfo {
	return s.build
}

// UpdateStatus represents update availability
type UpdateStatus string

const (
	StatusUpToDate        UpdateStatus = "up_to_date"
	StatusUpdateAvailable UpdateStatus = "update_available"
	StatusUnknown         UpdateStatus = "unknown"
)

// UpdateInfo contains update check results
type UpdateInfo struct {
	Status  UpdateStatus        `json:"status"`
	Reason  string              `json:"reason,omitempty"`
	Current CurrentVersionInfo  `json:"current"`
	Latest  *LatestVersionInfo  `json:"latest,omitempty"`
}

type CurrentVersionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
}

type LatestVersionInfo struct {
	Version     string `json:"version"`
	Tag         string `json:"tag"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt,omitempty"`
	Notes       string `json:"notes,omitempty"`
	Commit      string `json:"commit,omitempty"`
	Ref         string `json:"ref,omitempty"`
}

// CheckUpdate checks for available updates
func (s *VersionService) CheckUpdate(ctx context.Context) (UpdateInfo, error) {
	info := UpdateInfo{
		Current: CurrentVersionInfo{
			Version: s.build.Version,
			Commit:  s.build.Commit,
		},
	}

	// Dev builds: always unknown
	if s.build.IsDev() {
		info.Status = StatusUnknown
		info.Reason = "dev_build"
		return info, nil
	}

	// Prerelease builds: always unknown
	if s.build.IsPrerelease() {
		info.Status = StatusUnknown
		info.Reason = "prerelease_build"
		return info, nil
	}

	// Edge builds: compare commits
	if s.build.IsEdge() {
		return s.checkEdgeUpdate(ctx, info)
	}

	// Stable releases: compare semver
	return s.checkStableUpdate(ctx, info)
}

func (s *VersionService) checkEdgeUpdate(ctx context.Context, info UpdateInfo) (UpdateInfo, error) {
	if info.Current.Commit == "" {
		info.Status = StatusUnknown
		info.Reason = "missing_commit"
		return info, nil
	}

	// Fetch latest commit on main with caching
	commit, err := s.getMainHeadCommit(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Failed to fetch GitHub main HEAD")
		info.Status = StatusUnknown
		info.Reason = "github_error"
		return info, nil
	}

	latestCommit := commit.SHA[:7] // Short SHA

	if info.Current.Commit == latestCommit {
		info.Status = StatusUpToDate
	} else {
		info.Status = StatusUpdateAvailable
		info.Latest = &LatestVersionInfo{
			Version:     "edge",
			Tag:         "edge",
			URL:         commit.HTMLURL,
			Commit:      latestCommit,
			Ref:         "main",
			PublishedAt: commit.Commit.Author.Date.Format(time.RFC3339),
		}
	}

	return info, nil
}

func (s *VersionService) checkStableUpdate(ctx context.Context, info UpdateInfo) (UpdateInfo, error) {
	currentVer, err := semver.Parse(s.build.Version)
	if err != nil {
		s.logger.Warn().Err(err).Str("version", s.build.Version).Msg("Failed to parse current version")
		info.Status = StatusUnknown
		info.Reason = "invalid_version"
		return info, nil
	}

	// Fetch latest release with caching
	release, err := s.getLatestRelease(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("Failed to fetch GitHub latest release")
		info.Status = StatusUnknown
		info.Reason = "github_error"
		return info, nil
	}

	latestVer, err := semver.Parse(release.TagName)
	if err != nil {
		s.logger.Warn().Err(err).Str("tag", release.TagName).Msg("Failed to parse latest release version")
		info.Status = StatusUnknown
		info.Reason = "invalid_latest_version"
		return info, nil
	}

	if currentVer.LessThan(latestVer) {
		info.Status = StatusUpdateAvailable
		info.Latest = &LatestVersionInfo{
			Version:     release.TagName,
			Tag:         release.TagName,
			URL:         release.HTMLURL,
			PublishedAt: release.PublishedAt.Format(time.RFC3339),
			Notes:       release.Body,
		}
	} else {
		info.Status = StatusUpToDate
	}

	return info, nil
}

// Cached GitHub API calls
func (s *VersionService) getLatestRelease(ctx context.Context) (*github.Release, error) {
	cacheKey := "github_latest_release"
	release, err := getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*github.Release, error) {
		return s.gh.GetLatestRelease(ctx)
	}, UpdateTTL)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

func (s *VersionService) getMainHeadCommit(ctx context.Context) (*github.Commit, error) {
	cacheKey := "github_main_head"
	commit, err := getOrFetchFromCache(ctx, s.repo, s.logger, cacheKey, func() (*github.Commit, error) {
		return s.gh.GetMainHeadCommit(ctx)
	}, UpdateTTL)
	if err != nil {
		return nil, err
	}
	return &commit, nil
}
