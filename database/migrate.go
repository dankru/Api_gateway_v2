package database

import (
	"database/sql"
	"embed"

	// Import for side effects - needed for initializing the PostgresSQL driver
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(url string) error {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return errors.Wrap(err, "cannot connect to db to migrate")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "cannot ping db")
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "cannot set migration dialect")
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return errors.Wrap(err, "cannot get migration version")
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		if err = goose.DownTo(db, "migrations", version); err != nil {
			log.Err(err).Msgf("cannot rollbback migrations to version: %d", version)
		}

		return errors.Wrap(err, "cannot up migrations")
	}

	return nil
}
