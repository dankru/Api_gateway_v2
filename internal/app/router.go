package app

import (
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handler handler.Handler

	App *fiber.App
}

func NewRouter(config fiber.Config, handler handler.Handler) *Router {
	app := fiber.New(config)
	return &Router{handler: handler, App: app}
}

func (r *Router) InitializeRoutes() {
	r.App.Get("/user/", r.handler.GetUser)
}
