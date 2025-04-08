package main

import (
	"github.com/dankru/Api_gateway_v2/internal/app"
	"github.com/dankru/Api_gateway_v2/internal/handler"
	"github.com/dankru/Api_gateway_v2/internal/storage"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	connection, err := storage.GetConnect("postgres://postgres:postgres@postgres:5432/postgres")
	if err != nil {
		log.Fatalf("Failed to acquire pool: %s", err.Error())
	}

	routes := handler.NewHandler(connection)
	router := app.NewRouter(fiber.Config{AppName: "api_gateway"}, *routes)

	router.InitializeRoutes()

	err = router.App.Listen(":8000")
	if err != nil {
		//TODO: Log
		log.Fatalf("Unable to serve: %s", err.Error())
	}
}
