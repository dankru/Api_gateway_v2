package main

import (
	"os"

	"github.com/dankru/Api_gateway_v2/internal/app"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Err(err).Msg("app run failed")
		os.Exit(1)
	}
}
