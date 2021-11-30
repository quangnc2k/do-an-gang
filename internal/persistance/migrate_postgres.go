package persistance

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/quangnc2k/do-an-gang/internal/config"
)

type PostgresMigrator struct{}

func (pm *PostgresMigrator) Migrate(ctx context.Context, target int) (err error) {
	e := config.Env
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		e.PostgresUsername, e.PostgresPassword, e.PostgresHost, e.PostgresPort, e.PostgresDb, "disable",
	)

	db, err := sql.Open("postgres", conn)

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;")
	if err != nil {
		return
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", e.MigrateFilePath), "postgres", driver)
	if err != nil {
		return
	}

	if target >= 0 {
		return m.Migrate(uint(target))
	}

	return m.Up()
}

func NewPostgresMigrator(ctx context.Context) *PostgresMigrator{
	return &PostgresMigrator{}
}
