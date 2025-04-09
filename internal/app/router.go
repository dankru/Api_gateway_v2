package app

import (
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Router struct {
	handler *handler.Handler

	App *fiber.App
}

func NewRouter(config fiber.Config, handler *handler.Handler) *Router {
	app := fiber.New(config)
	router := &Router{handler: handler, App: app}
	log.Info().Msg("Initializing routes")
	router.InitializeRoutes()
	return router
}

func (r *Router) InitializeRoutes() {
	r.App.Get("/user/:id", r.handler.GetUser)
	r.App.Put("/user/:id", r.handler.ReplaceUser)
	r.App.Delete("/user/:id", r.handler.DeleteUser)
	r.App.Post("/user", r.handler.CreateUser)

	routes := r.App.GetRoutes()
	for _, route := range routes {
		log.Info().Msgf("Initialized route: %s [%s]", route.Path, route.Method)
	}
}
