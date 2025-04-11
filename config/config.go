package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type config struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	AppName     string
	Address     string
}

func Init() (config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
		return config{}, err
	}

	viper.AutomaticEnv()

	config := config{
		DB_USER:     viper.GetString("DB_USER"),
		DB_PASSWORD: viper.GetString("DB_PASSWORD"),
		DB_HOST:     viper.GetString("DB_HOST"),
		DB_PORT:     viper.GetString("DB_PORT"),
		DB_NAME:     viper.GetString("DB_NAME"),
		AppName:     viper.GetString("app.name"),
		Address:     viper.GetString("app.port"),
	}

	return config, nil
}
