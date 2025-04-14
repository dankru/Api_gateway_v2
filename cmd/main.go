package main

import (
	"github.com/dankru/Api_gateway_v2/internal/app"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		log.Err(err).Msg("app run failed")
		os.Exit(1)
	}
}
