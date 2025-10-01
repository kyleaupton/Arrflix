package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(env string) zerolog.Logger {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if env == "dev" {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	}
	return l
}
