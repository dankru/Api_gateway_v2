package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

func Init() error {
	logLevel := viper.GetString("log.level")

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Err(err).Msg("logger global level parsing failed")
		level = zerolog.InfoLevel
		return err
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("global log level is set to: %s", zerolog.GlobalLevel())
	return nil
}
