package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

type DB struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type Cache struct {
	TTL             time.Duration
	CleanerInterval time.Duration
}

type Log struct {
	Level string
}

type Metrics struct {
	Port         string
	SendInterval time.Duration
}

type App struct {
	Name    string
	Address string
	Cache   Cache
	Log     Log
	Metrics Metrics
}

type Config struct {
	DB
	App
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

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("unable to decode config into struct")
		return Config{}, err
	}
	fmt.Println("All settings:", viper.AllSettings()) // показать все загруженные ключи
	fmt.Println("CONFIG: ", cfg)
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
