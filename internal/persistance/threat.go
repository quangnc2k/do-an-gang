package persistance

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"time"
)

type ThreatRepository interface {
	Paginate(ctx context.Context, page, perPage int, orderBy, search string, start, end time.Time) ([][]interface{}, error)
	StatsByPhase(ctx context.Context, start, end time.Time) (map[string]int, error)
	StatsBySeverity(ctx context.Context, start, end time.Time) (map[string]int, error)
	TopHostAffected(ctx context.Context, start, end time.Time) (map[string]int64, error)
	TopAttacker(ctx context.Context, start, end time.Time) (map[string]int64, error)
	StoreThreatInBatch(ctx context.Context, threatChan []model.Threat) (err error)
}
