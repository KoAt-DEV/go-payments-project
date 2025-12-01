package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Options struct {
	Environment string
	ServiceName string
}

func New(opts Options) *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	if opts.Environment != "production" {
		console := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}

		log.Logger = log.Output(console).Level(zerolog.DebugLevel)
	} else {
		log.Logger = log.Output(os.Stdout).Level(zerolog.InfoLevel)
	}

	logger := log.With().
		Str("service", opts.ServiceName).
		Timestamp().
		Caller().
		Logger()

	return &logger
}
