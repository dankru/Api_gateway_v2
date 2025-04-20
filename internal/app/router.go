package app

import (
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/dankru/Api_gateway_v2/internal/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func newRouter(config fiber.Config, handler *handler.Handler) *fiber.App {
	app := fiber.New(config)
	log.Info().Msg("Initializing routes")
	user := app.Group("/user")

	user.Use(metrics.PrometheusMiddleware())

	user.Get("/:id", handler.GetUser)
	user.Put("/:id", handler.ReplaceUser)
	user.Post("/", handler.CreateUser)
	user.Delete("/:id", handler.DeleteUser)

	routes := app.GetRoutes()
	for _, route := range routes {
		log.Info().Msgf("Initialized route: %s [%s]", route.Path, route.Method)
	}

	return app
}
