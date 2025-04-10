package logger

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

func LoggerInit() {
	logLevel := viper.GetString("log.level")

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Err(errors.Wrap(err, "failed to parse global logging level from config")).Msg("logger global level parsing failed")
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("global log level is set to: %s", zerolog.GlobalLevel())
}
