package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

type DB struct {
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
}

type Cache struct {
	TTL             time.Duration `mapstructure:"ttl"`
	CleanerInterval time.Duration `mapstructure:"cleanerInterval"`
}

type Log struct {
	Level string `mapstructure:"level"`
}

type Metrics struct {
	Port string `mpstructure:"port"`
}

type App struct {
	Name    string  `mapstructure:"name"`
	Address string  `mapstructure:"port"`
	Cache   Cache   `mapstructure:"cache"`
	Log     Log     `mapstructure:"log"`
	Metrics Metrics `mapstructure:"metrics"`
}

type Config struct {
	DB  `mapstructure:",squash"`
	App `mapstructure:"app"`
}

func Init() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
		return Config{}, err
	}

	viper.AutomaticEnv()

	envKeys := []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	}

	for _, key := range envKeys {
		if err := viper.BindEnv(key); err != nil {
			log.Warn().Err(err).Msgf("viper failed to bind key from env: %s", key)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("unable to decode config into struct")
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name)
}
