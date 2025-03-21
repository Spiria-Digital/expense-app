package migrations

import (
	"context"
	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/models"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().
			Model((*models.Expense)(nil)).
			ForeignKey("(owner_id) REFERENCES users (id) ON DELETE CASCADE").
			ForeignKey("(category_id) REFERENCES categories (id) ON DELETE SET NULL").
			IfNotExists().
			Exec(ctx)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().
			Model((*models.Expense)(nil)).
			IfExists().
			Exec(ctx)
		return err
	})
}
