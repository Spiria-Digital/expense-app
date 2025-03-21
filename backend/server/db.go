package server

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/storage"
)

var BunDB *bun.DB

func init() {
	filePath := filepath.Join(storage.GetRootDir(), "expenses.db")
	db, err := storage.NewBunDB(filePath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create bun db")
	}
	BunDB = db
}
