package services

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
)

func Migrate(ctx context.Context) (err error){
	return persistance.GetMigrator().Migrate(ctx, -1)
}
