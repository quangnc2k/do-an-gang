package persistance

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"time"
)

type AlertRepository interface {
	FindAll(ctx context.Context, showAll bool) (alerts []model.Alert, err error)
	Create(ctx context.Context, alert model.Alert) (err error)
	Count(ctx context.Context, from, to time.Time) (count int, err error)
	Resolve(ctx context.Context, resolved bool, resolvedBy string, ids ...string) (rows int, err error)
}
