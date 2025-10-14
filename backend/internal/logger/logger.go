package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func New(dev bool) *Logger {
	if dev {
		cw := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05.000", // hh:mm:ss.mmm
		}

		l := zerolog.New(cw).
			With().
			Timestamp().
			Caller().
			Logger()

		// Prefer a readable time in dev
		zerolog.TimeFieldFormat = time.RFC3339Nano
		return &l
	}

	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &l
}
