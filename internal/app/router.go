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
	r.App.Get("/user/:id", r.handler.GetUser)
	r.App.Put("/user/:id", r.handler.ReplaceUser)
	r.App.Delete("/user/:id", r.handler.DeleteUser)
	r.App.Post("/user", r.handler.CreateUser)
}
