package persistance

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"time"
)

type ThreatRepository interface {
	Paginate(ctx context.Context, page, perPage int, orderBy, search string, start, end time.Time) (model.PaginateData, error)
	StatsByPhase(ctx context.Context, start, end time.Time) (model.PieChartData, error)
	StatsBySeverity(ctx context.Context, start, end time.Time) (model.PieChartData, error)
	TopHostAffected(ctx context.Context, start, end time.Time) (model.BarChartData, error)
	TopAttacker(ctx context.Context, start, end time.Time) (model.BarChartData, error)
	HistogramAffected(ctx context.Context, width string, start, end time.Time) (model.LineChartData, error)
	StoreThreatInBatch(ctx context.Context, threatChan []model.Threat) (err error)

	RecentAffected(ctx context.Context) (host map[string]time.Time, err error)
	RecentAttackByPhase(ctx context.Context) (host map[string]time.Time, err error)
	Overview(ctx context.Context) (total, recent, numOfHost int64, err error)
}
