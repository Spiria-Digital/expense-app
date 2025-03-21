package main

import (
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"

	"github.com/Spiria-Digital/expense-manager/cmd/bun/migrations"
	"github.com/Spiria-Digital/expense-manager/server"
)

func main() {
	app := &cli.App{
		Name: "bun",
		Commands: []*cli.Command{
			subCommands(migrate.NewMigrator(server.BunDB, migrations.Migrations)),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Error running database migrations")
	}
}

func subCommands(migrator *migrate.Migrator) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "initialize database",
				Action: func(c *cli.Context) error {
					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context)

					group, err := migrator.Migrate(c.Context)
					if err != nil {
						return err
					}

					if group.IsZero() {
						log.Info().Msg("no new migrations to run")
					} else {
						log.Info().Msgf("migrations applied: (%s)", group)
					}

					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback database",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context)

					group, err := migrator.Rollback(c.Context)
					if err != nil {
						return err
					}

					if group.IsZero() {
						log.Info().Msg("no migrations to rollback")
					} else {
						log.Warn().Msgf("migrations rolled back: (%s)", group)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "show migration status",
				Action: func(c *cli.Context) error {
					status, err := migrator.MigrationsWithStatus(c.Context)
					if err != nil {
						return err
					}

					log.Info().
						Any("status", status).
						Any("notApplied", status.Unapplied()).
						Any("applied", status.Applied()).
						Msg("Migration status")

					return nil
				},
			},
			{
				Name:  "create",
				Usage: "create new migration",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					if name == "" {
						return errors.New("missing migration name")
					}
					mf, err := migrator.CreateGoMigration(c.Context, name)
					if err != nil {
						return err
					}
					log.Info().Str("file", mf.Path).Msg("Migration created")
					return nil
				},
			},
		},
	}
}
