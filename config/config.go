package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func ConfigInit() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
	}
	viper.AutomaticEnv()
}
