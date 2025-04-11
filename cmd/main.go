package main

import (
	"context"
	"fmt"
	"github.com/dankru/Api_gateway_v2/cache"
	"github.com/dankru/Api_gateway_v2/config"
	"github.com/dankru/Api_gateway_v2/internal/app"
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/dankru/Api_gateway_v2/internal/storage"
	"github.com/dankru/Api_gateway_v2/internal/usecase"
	"github.com/dankru/Api_gateway_v2/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Msg("failed to initialize configs")
	}

	if err := logger.Init(); err != nil {
		log.Error().Msg("failed to initialize logger")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_NAME)

	log.Info().Msgf("initializing db connection: %s", connStr)
	conn, err := storage.GetConnect(connStr)
	if err != nil {
		log.Fatal().Err(err).
			Msg("failed to get db pool")
	}

	repo := repository.NewUserRepository(conn)
	cacheDecorator := cache.NewCacheDecorator(repo, cfg)
	cacheDecorator.StartCleaner(context.Background())
	uc := usecase.NewUseCase(cacheDecorator)
	handle := handler.NewHandler(uc)

	router := app.NewRouter(fiber.Config{AppName: cfg.AppName}, handle)

	log.Info().Msgf("listen and serve on: %s", cfg.Address)
	if err := router.Listen(cfg.Address); err != nil {
		log.Fatal().
			Err(err).
			Msgf("unable to listen and serve on %s", cfg.Address)
	}

}
