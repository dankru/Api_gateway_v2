package logger

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func Init(logLevel string) error {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Err(err).Msg("logger global level parsing failed")
		level = zerolog.InfoLevel
		return errors.Wrap(err, "failed to parse log level from conf")
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("global log level is set to: %s", zerolog.GlobalLevel())
	return nil
}
