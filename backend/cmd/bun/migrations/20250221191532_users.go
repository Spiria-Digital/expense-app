package migrations

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/models"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().
			Model((*models.User)(nil)).
			IfNotExists().
			Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().
			Model((*models.User)(nil)).
			IfExists().
			Exec(ctx)
		return err
	})
}
