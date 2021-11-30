package migrate

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/services"
)

func Migrate(ctx context.Context) (err error) {
	return services.Migrate(ctx)
}
