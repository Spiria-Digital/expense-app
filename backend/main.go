package main

import (
	"github.com/rs/zerolog/log"

	"github.com/Spiria-Digital/expense-manager/server"
)

func main() {
	if err := server.Router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal().Err(err).Msg("error running app")
	}
}
