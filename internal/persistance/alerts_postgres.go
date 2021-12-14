package persistance

import (
	"context"
	"encoding/json"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"time"
)

type AlertRepositorySQL struct {
	connection *pgxpool.Pool
}

func NewAlertSQLRepo(ctx context.Context, conn *pgxpool.Pool) (AlertRepository, error) {
	if conn == nil {
		return nil, errors.New("invalid sql connection")
	}
	return &AlertRepositorySQL{connection: conn}, nil
}

func (r *AlertRepositorySQL) FindAll(ctx context.Context, showAll bool) (alerts []model.Alert, err error) {
	var query string

	if showAll {
		query = `SELECT id, created_at, details, resolved_at, resolved_by
				FROM alerts`
	} else {
		query = `SELECT id, created_at, details, resolved_at, resolved_by
				FROM alerts
				WHERE resolved = FALSE`
	}

	rows, err := r.connection.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var alert model.Alert
		err = rows.Scan(&alert.ID, alert.CreatedAt, alert.Details, alert.ResolvedAt, alert.ResolvedBy)
		if err != nil {
			continue
		}

		alerts = append(alerts, alert)
	}

	return
}

func (r *AlertRepositorySQL) Create(ctx context.Context, alert model.Alert) (err error) {
	d, err := json.Marshal(alert.Details)
	if err != nil {
		return
	}
	query := `INSERT INTO alerts (id, created_at, details, resolved) VALUES ($1, $2, $3, $4)`
	_, err = r.connection.Exec(ctx, query, alert.ID, alert.CreatedAt, d, alert.Resolved)
	return
}

func (r *AlertRepositorySQL) Count(ctx context.Context, from, to time.Time) (count int, err error) {
	err = r.connection.QueryRow(ctx, "SELECT count(id) FROM alerts WHERE created_at >= $1 AND created_at <= $2", from, to).Scan(&count)
	return
}

func (r *AlertRepositorySQL) Resolve(ctx context.Context, resolved bool, resolvedBy string, ids ...string) (rows int, err error) {
	if len(ids) == 0 {
		return 0, nil
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	builder := psql.Update("alerts")
	if resolved {
		builder = builder.SetMap(map[string]interface{}{
			"resolved_at": time.Now(),
			"resolved_by": resolvedBy,
		})
	} else {
		builder = builder.SetMap(map[string]interface{}{
			"resolved_at": nil,
			"resolved_by": nil,
		})
	}
	ors := sq.Or{}
	for _, id := range ids {
		ors = append(ors, sq.Eq{"id": id})
	}
	builder = builder.Where(ors)
	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	tag, err := r.connection.Exec(ctx, sql, args...)
	if err != nil {
		return
	}
	rows = int(tag.RowsAffected())
	return
}
