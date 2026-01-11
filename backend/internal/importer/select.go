package importer

import (
	"path/filepath"
	"strings"

	"github.com/kyleaupton/Arrflix/internal/downloader"
)

// PickMainMovieFile chooses the "main" file for a movie download.
// Strategy: largest video file by extension, excluding obvious samples.
func PickMainMovieFile(files []downloader.File) (downloader.File, bool) {
	var (
		best    downloader.File
		bestSet bool
	)
	for _, f := range files {
		if f.Size <= 0 {
			continue
		}
		if !IsVideoPath(f.Path) {
			continue
		}
		if LooksLikeSample(f.Path) {
			continue
		}
		if !bestSet || f.Size > best.Size {
			best = f
			bestSet = true
		}
	}
	if bestSet {
		return best, true
	}

	// Fallback: any largest file (still ignore samples if possible)
	for _, f := range files {
		if f.Size <= 0 {
			continue
		}
		if LooksLikeSample(f.Path) {
			continue
		}
		if !bestSet || f.Size > best.Size {
			best = f
			bestSet = true
		}
	}
	if bestSet {
		return best, true
	}

	// Last resort: any file
	for _, f := range files {
		if !bestSet || f.Size > best.Size {
			best = f
			bestSet = true
		}
	}
	return best, bestSet
}

func EnsureExt(path, ext string) string {
	if ext == "" {
		return path
	}
	if strings.EqualFold(filepath.Ext(path), ext) {
		return path
	}
	return path + ext
}
