package persistance

import "context"

var migrator Migrator

type Migrator interface {
	Migrate(ctx context.Context, target int) (err error)
}

func GetMigrator() Migrator{
	return migrator
}

func LoadMigratorWithPgx(ctx context.Context) {
	migrator = NewPostgresMigrator(ctx)
}