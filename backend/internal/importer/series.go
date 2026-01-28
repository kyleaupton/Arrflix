package importer

import (
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/kyleaupton/arrflix/internal/downloader"
)

// SeriesInfo contains parsed season and episode numbers.
type SeriesInfo struct {
	Season   int
	Episodes []int
}

var (
	// Standard S01E01 pattern
	reStandard = regexp.MustCompile(`(?i)s(\d+)e(\d+)(?:-?e?(\d+))?`)
	// 1x01 pattern
	reAlternative = regexp.MustCompile(`(?i)(\d+)x(\d+)`)
)

// ParseSeriesInfo extracts season and episode information from a filename.
func ParseSeriesInfo(filename string) (SeriesInfo, bool) {
	// Try standard S01E01
	if matches := reStandard.FindStringSubmatch(filename); len(matches) >= 3 {
		season, _ := strconv.Atoi(matches[1])
		ep1, _ := strconv.Atoi(matches[2])
		info := SeriesInfo{Season: season, Episodes: []int{ep1}}

		// Multi-episode S01E01E02 or S01E01-E02 or S01E01-02
		if len(matches) >= 4 && matches[3] != "" {
			ep2, _ := strconv.Atoi(matches[3])
			// Add all episodes in range if it looks like a range, or just the second one
			// For now, let's just add the second one.
			if ep2 > ep1 && ep2-ep1 < 10 { // sanity check on range
				for i := ep1 + 1; i <= ep2; i++ {
					info.Episodes = append(info.Episodes, i)
				}
			} else {
				info.Episodes = append(info.Episodes, ep2)
			}
		}
		return info, true
	}

	// Try 1x01
	if matches := reAlternative.FindStringSubmatch(filename); len(matches) == 3 {
		season, _ := strconv.Atoi(matches[1])
		ep, _ := strconv.Atoi(matches[2])
		return SeriesInfo{Season: season, Episodes: []int{ep}}, true
	}

	return SeriesInfo{}, false
}

// MatchFilesToEpisodes matches downloader files to their corresponding episodes.
func MatchFilesToEpisodes(files []downloader.File, targetSeason *int, targetEpisode *int) map[int]downloader.File {
	matched := make(map[int]downloader.File)

	for _, f := range files {
		if !IsVideoPath(f.Path) || LooksLikeSample(f.Path) {
			continue
		}

		filename := filepath.Base(f.Path)
		info, ok := ParseSeriesInfo(filename)
		if !ok {
			// If we are looking for a specific episode and this is the only large video file,
			// it might be a renamed episode file.
			continue
		}

		// If a target season is specified, it must match.
		if targetSeason != nil && info.Season != *targetSeason {
			continue
		}

		// If a target episode is specified, it must be in the parsed episodes.
		if targetEpisode != nil {
			found := false
			for _, ep := range info.Episodes {
				if ep == *targetEpisode {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Map each episode found in the file to this file.
		for _, ep := range info.Episodes {
			// If multiple files match the same episode, keep the largest one.
			if existing, ok := matched[ep]; !ok || f.Size > existing.Size {
				matched[ep] = f
			}
		}
	}

	return matched
}
