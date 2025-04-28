package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dankru/Api_gateway_v2/config"
	"github.com/dankru/Api_gateway_v2/database"
	"github.com/dankru/Api_gateway_v2/internal/cache"
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/dankru/Api_gateway_v2/internal/metrics"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/dankru/Api_gateway_v2/internal/storage"
	"github.com/dankru/Api_gateway_v2/internal/tracing"
	"github.com/dankru/Api_gateway_v2/internal/usecase"
	"github.com/dankru/Api_gateway_v2/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Init()
	if err != nil {
		log.Error().Msg("failed to initialize configs")
		return errors.Wrap(err, "config initilization failed")
	}
	if err := logger.Init(cfg.Log.Level); err != nil {
		log.Error().Msg("failed to initialize logger")
		return errors.Wrap(err, "logger initialization failed")
	}

	connStr := cfg.GetConnStr()
	if err := database.Migrate(connStr); err != nil {
		log.Err(err).Msg("failed to migrate")
	}

	log.Info().Msgf("initializing db connection: %s", connStr)
	conn, err := storage.GetConnect(connStr)
	if err != nil {
		log.Error().Err(err).
			Msg("failed to get db pool")
		return errors.Wrap(err, "initializing db connection failed")
	}
	defer conn.Close()

	tracerProvider, err := tracing.NewTracerProvider(cfg.App.Name, cfg.Jaeger.Collector.Endpoint, cfg.App.Environment)
	if err != nil {
		log.Error().Msg("failed to initialize jaeger")
		return errors.Wrap(err, "failed to initialize jaeger")
	}
	otel.SetTracerProvider(tracerProvider)

	repo := repository.NewUserRepository(conn)
	cacheDecorator := cache.NewCacheDecorator(repo, cfg.App.Cache.TTL)
	uc := usecase.NewUserUsecase(cacheDecorator)
	handle := handler.NewHandler(uc)
	//
	cacheDecorator.StartCleaner(ctx, cfg.App.Cache.CleanerInterval)

	metrics.InitMetrics(cfg.App.Metrics.Port, cacheDecorator, cfg.Metrics.SendInterval)

	router := newRouter(fiber.Config{AppName: cfg.App.Name}, handle)
	go func() {
		log.Info().Msgf("listen and serve on: %s", cfg.App.Address)
		if err := router.Listen(":" + cfg.App.Address); err != nil {
			log.Error().
				Err(err).
				Msgf("unable to listen and serve on %s", cfg.App.Address)
		}
	}()

	<-stop
	log.Info().Msg("shutting down gracefully")

	cancel()
	if err := router.Shutdown(); err != nil {
		log.Error().Err(err).Msg("error shutting down server")
		return errors.Wrap(err, "server shutdown failed")
	}

	log.Info().Msg("server stopped gracefully")

	return nil
}
