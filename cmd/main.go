package main

import (
	"github.com/dankru/Api_gateway_v2/internal/app"
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/dankru/Api_gateway_v2/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var port = ":8000"
var connString = "postgres://postgres:postgres@postgres:5432/postgres"

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msgf("Initializing db connection: %s", connString)
	connection, err := storage.GetConnect(connString)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get db pool")
	}

	routes := handler.NewHandler(connection)
	router := app.NewRouter(fiber.Config{AppName: "api_gateway"}, *routes)

	log.Info().Msg("Initializing routes")
	router.InitializeRoutes()

	log.Info().Msgf("Listen and serve on: %s", port)
	err = router.App.Listen(port)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to listen and serve on %s", port)
	}
}
