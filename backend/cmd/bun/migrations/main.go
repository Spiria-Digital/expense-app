package migrations

import (
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
)

var (
	Migrations = migrate.NewMigrations()
)

func init() {
	if err := Migrations.DiscoverCaller(); err != nil {
		log.Fatal().Err(err).Msg("failed to discover caller")
	}
}
