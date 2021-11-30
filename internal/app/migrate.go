package app

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/app/migrate"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
)

func Migrate(ctx context.Context) (err error) {
	config.Init()

	persistance.LoadRepoContainerWithPgx(ctx)
	persistance.LoadMigratorWithPgx(ctx)

	return migrate.Migrate(ctx)
}
