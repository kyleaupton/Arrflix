package tmdb

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/kyleaupton/snaggle/backend/internal/config"
)

var Client *tmdb.Client

func InitTmdb() {
	client, err := tmdb.Init(config.Load().TmdbAPIKey)
	if err != nil {
		panic(err)
	}

	Client = client
}
