package main

import (
	"fmt"
	"github.com/dankru/Api_gateway_v2/cmd/logger"
	"github.com/dankru/Api_gateway_v2/config"
	"github.com/dankru/Api_gateway_v2/internal/app"
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/dankru/Api_gateway_v2/internal/repository"
	"github.com/dankru/Api_gateway_v2/internal/storage"
	"github.com/dankru/Api_gateway_v2/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {

	config.ConfigInit()

	logger.LoggerInit()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"))

	log.Info().Msgf("initializing db connection: %s", connStr)
	conn, err := storage.GetConnect(connStr)
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, "connect to db failed")).
			Msg("failed to get db pool")
	}

	repo := repository.NewUserRepository(conn)
	uc := usecase.NewUseCase(repo)
	handle := handler.NewHandler(uc)

	router := app.NewRouter(fiber.Config{AppName: "api_gateway"}, handle)

	log.Info().Msgf("listen and serve on: %s", viper.GetString("app.port"))
	if err := router.Listen(viper.GetString("app.port")); err != nil {
		log.Fatal().
			Err(errors.Wrap(err, "failed to start server")).
			Msgf("unable to listen and serve on %s", viper.GetString("app.port"))
	}
}
