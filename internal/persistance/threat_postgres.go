package persistance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quangnc2k/do-an-gang/internal/model"
)

type ThreatRepositorySQL struct {
	connection *pgxpool.Pool
}

func (r *ThreatRepositorySQL) Paginate(ctx context.Context, page, perPage int, orderBy, search string, start, end time.Time) (output [][]interface{}, err error) {
	header := []interface{}{"created_at", "seen_at", "src_host", "dst_host", "confidence", "severity", "phase"}
	output = append(output, header)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("created_at", "seen_at", "src_host", "dst_host", "confidence", "severity", "phase")
	queryBuilder = queryBuilder.From("threats")
	queryBuilder = queryBuilder.Where(
		sq.And{
			sq.GtOrEq{"created_at": "?"},
			sq.LtOrEq{"created_at": "?"},
			sq.Or{
				sq.Like{"src_host": "?"},
				sq.Like{"dst_host": "?"},
				sq.Like{"phase": fmt.Sprintf("%%%s%%", "?")},
			},
		}, start, end, search, search, search)
	queryBuilder = queryBuilder.Offset(uint64(int64(page * perPage)))
	queryBuilder = queryBuilder.Limit(uint64(int64(perPage)))
	queryBuilder.OrderBy(orderBy)

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return
	}

	rows, err := r.connection.Query(ctx, query)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		row := make([]interface{}, 7)

		err = rows.Scan(&row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6])
		if err != nil {
			log.Println(err)
			continue
		}

		output = append(output, row)
	}

	return
}

func (r *ThreatRepositorySQL) StatsByPhase(ctx context.Context, start, end time.Time) (map[string]int, error) {
	var output map[string]int

	query := `SELECT
    		phase, SUM(count)
			FROM threats
			WHERE _bucket >= $1 AND _bucket <= $2
			GROUP BY phase
			ORDER BY phase`

	rows, err := r.connection.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var phase string
		var total int

		err = rows.Scan(&phase, &total)
		if err != nil {
			return nil, err
		}

		output[phase] = total
	}

	return output, err
}

func (r *ThreatRepositorySQL) StatsBySeverity(ctx context.Context, start, end time.Time) (map[string]int, error) {
	var output map[string]int

	query := `SELECT
    		severity, SUM(count)
			FROM stats_threat_by_severity
			WHERE _bucket >= $1 AND _bucket <= $2
			GROUP BY severity
			ORDER BY severity`

	rows, err := r.connection.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var severity string
		var total int

		err = rows.Scan(&severity, &total)
		if err != nil {
			return nil, err
		}

		output[severity] = total
	}

	return output, err
}

func (r *ThreatRepositorySQL) TopHostAttacked(ctx context.Context, start, end time.Time) (map[string]int64, error) {
	panic("implement me")
}

func (r *ThreatRepositorySQL) TopAttacker(ctx context.Context, start, end time.Time) (map[string]int64, error) {
	panic("implement me")
}

func NewThreatSQLRepo(ctx context.Context, conn *pgxpool.Pool) (ThreatRepository, error) {
	if conn == nil {
		return nil, errors.New("invalid sql connection")
	}
	return &ThreatRepositorySQL{connection: conn}, nil
}

func (r *ThreatRepositorySQL) StoreThreatInBatch(ctx context.Context, threats []model.Threat) (err error) {
	batch := &pgx.Batch{}

	if len(threats) == 0 {
		return
	}

	for _, threat := range threats {
		query := "INSERT INTO threats (id, created_at, seen_at, src_host, dst_host, conn_id, severity, confidence, phase, metadata)" +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
		batch.Queue(query,
			threat.ID,
			time.Now(),
			threat.SeenAt,
			threat.SourceHost,
			threat.DestinationHost,
			threat.ConnID,
			threat.Severity,
			threat.Confidence,
			threat.Phase,
			threat.Metadata,
		)

	}

	batchResult := r.connection.SendBatch(ctx, batch)

	defer batchResult.Close()

	_, err = batchResult.Exec()
	if err != nil {
		return
	}

	return nil
}
