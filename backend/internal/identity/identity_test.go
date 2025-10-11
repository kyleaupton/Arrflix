package identity

import (
	"testing"
)

func TestGetIdFromString(t *testing.T) {
	type Test struct {
		Path             string
		ExpectedProvider string
		ExpectedID       string
	}

	tests := []Test{
		{
			Path:             "/mnt/media-pipeline/media/tv/South Park (1997) {tvdb-75897}/Season 27/South Park (1997) - S27E02 - Got A Nut [WEBDL-1080p][EAC3 5.1][h264]-FLUX.mkv",
			ExpectedProvider: "tvdb",
			ExpectedID:       "75897",
		},
		{
			Path:             "/mnt/media-pipeline/media/movies/A Minecraft Movie (2025) {tmdb-950387}.mp4",
			ExpectedProvider: "tmdb",
			ExpectedID:       "950387",
		},
		{
			Path:             "/mnt/media-pipeline/media/movies/A Minecraft Movie (2025) {imdb-tt3566834}.mp4",
			ExpectedProvider: "imdb",
			ExpectedID:       "tt3566834",
		},
		{
			Path:             "/media/movies/21 Jump Street (2012) {tmdb-64688}/21 Jump Street (2012) {tmdb-64688} [Bluray-1080p][DTS 5.1][x264]-CtrlHD.mkv",
			ExpectedProvider: "tmdb",
			ExpectedID:       "64688",
		},
	}

	for _, test := range tests {
		provider, id, err := getIdFromString(test.Path)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if provider != test.ExpectedProvider || id != test.ExpectedID {
			t.Errorf("Expected tmdb id %s, got %s", test.ExpectedID, id)
		}
	}
}

func TestGetSeasonAndEpisodeFromPath(t *testing.T) {
	type Test struct {
		Path            string
		ExpectedSeason  int32
		ExpectedEpisode int32
	}

	tests := []Test{
		{
			Path:            "/media/tv/Show Name (2025) {tmdb-123456}/Season 1/Show Name (2025) s01e01.mkv",
			ExpectedSeason:  1,
			ExpectedEpisode: 1,
		},
		{
			Path:            "/mnt/media-pipeline/media/tv/South Park (1997) {tvdb-75897}/Season 27/South Park (1997) - S27E02 - Got A Nut [WEBDL-1080p][EAC3 5.1][h264]-FLUX.mkv",
			ExpectedSeason:  27,
			ExpectedEpisode: 2,
		},
	}

	for _, test := range tests {
		season, episode, err := getSeasonAndEpisodeFromPath(test.Path)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if season != test.ExpectedSeason || episode != test.ExpectedEpisode {
			t.Errorf("Expected season %d and episode %d, got %d and %d", test.ExpectedSeason, test.ExpectedEpisode, season, episode)
		}
	}
}
