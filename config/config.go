package config

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

type Agent struct {
	Host string
	Port string
}

type Sampler struct {
	Type  string
	Param string
}

type Collector struct {
	Endpoint string
}

type Jaeger struct {
	Agent     Agent
	Collector Collector
	Sampler   Sampler
}

type App struct {
	Name        string
	Address     string
	Environment string
	Cache       Cache
	Log         Log
	Metrics     Metrics
}

type Config struct {
	DB
	App
	Jaeger
}

var AppName string

func Init() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
		return nil, err
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("unable to decode config into struct")
		return nil, err
	}

	setAppName(cfg.App.Name)
	return &cfg, nil
}

func (c *Config) GetConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name)
}

func setAppName(name string) {
	AppName = name
}
