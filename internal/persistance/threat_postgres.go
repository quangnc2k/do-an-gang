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

func (r *ThreatRepositorySQL) HistogramAffected(ctx context.Context, width string, start, end time.Time) (output model.LineChartData, err error) {
	query1 := `SELECT affected_host FROM stats_threat_by_affected
						WHERE t >= $1 AND t <= $2
						GROUP BY affected_host 
						ORDER BY SUM(count) DESC
						LIMIT 10`

	var affectedHosts []string

	rows, err := r.connection.Query(ctx, query1, start, end)
	if err != nil {
		return
	}

	for rows.Next() {
		var host string
		err = rows.Scan(&host)
		if err != nil {
			log.Println(err)
			continue
		}

		affectedHosts = append(affectedHosts, host)
	}

	for i, host := range affectedHosts {
		query := `SELECT time_bucket_gapfill($1, seen_at) as t, SUM(count) as total
					FROM stats_threat_by_affected
					WHERE t >= $2 AND t <= $3 AND affected_host = $4
					GROUP BY t, affected_host
					ORDER BY t ASC`

		rows, err := r.connection.Query(ctx, query, width, start, end, host)
		if err != nil {
			return output, err
		}

		var dataset model.Dataset

		for rows.Next() {
			var t time.Time
			var total int

			err := rows.Scan(&t, &total)
			if err != nil {
				rows.Close()
				return output, err
			}

			if i == 0 {
				output.Labels = append(output.Labels, t.String())
			}

			dataset.Data = append(dataset.Data, total)
		}

		dataset.Name = host
		output.Datasets = append(output.Datasets, dataset)
		rows.Close()
	}

	return
}

func (r *ThreatRepositorySQL) Paginate(ctx context.Context, page, perPage int, orderBy, search string, start, end time.Time) (output model.PaginateData, err error) {
	//header := []interface{}{"created_at", "seen_at", "affected_host", "attacker_host", "confidence", "severity", "phase"}
	//output = append(output, header)

	var count int
	queryCount := `SELECT COUNT(*) from threats WHERE created_at >= $1 AND created_at <= $2`
	err = r.connection.QueryRow(ctx, queryCount, start, end).Scan(&count)
	if err != nil {
		return
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	queryBuilder := psql.Select("id", "created_at", "seen_at", "affected_host", "attacker_host", "confidence", "severity", "phase", "metadata")
	queryBuilder = queryBuilder.From("threats")

	where := sq.And{
		sq.GtOrEq{"created_at": start},
		sq.LtOrEq{"created_at": end},
	}

	if search != "" {
		where = append(where, sq.Or{
			sq.Like{"affected_host": search},
			sq.Like{"attacker_host": search},
			sq.Like{"phase": fmt.Sprintf("%%%s%%", search)},
		})
	}

	queryBuilder = queryBuilder.Where(where)
	queryBuilder = queryBuilder.Offset(uint64(int64((page - 1) * perPage)))
	queryBuilder = queryBuilder.Limit(uint64(int64(perPage)))
	queryBuilder.OrderBy(orderBy)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return
	}

	rows, err := r.connection.Query(ctx, query, args...)
	if err != nil {
		return
	}

	output.Links.Pagination.Total = count
	output.Links.Pagination.PerPage = perPage
	output.Links.Pagination.CurrentPage = page
	output.Links.Pagination.From = (page-1)*perPage + 1
	output.Links.Pagination.To = page * perPage

	defer rows.Close()
	for rows.Next() {
		//row := make([]interface{}, 7)
		//
		//err = rows.Scan(&row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6])
		var createdAt, seenAt time.Time
		var id, affectedHost, attackHost, phase string
		var confidence float64
		var severity int
		var metadata map[string]interface{}

		err = rows.Scan(&id, &createdAt, &seenAt, &affectedHost, &attackHost, &confidence, &severity, &phase, &metadata)
		if err != nil {
			log.Println(err)
			continue
		}

		output.Data = append(output.Data, model.Threat{
			ID:           id,
			CreatedAt:    createdAt,
			SeenAt:       seenAt,
			AffectedHost: affectedHost,
			AttackerHost: attackHost,
			Confidence:   confidence,
			Severity:     severity,
			Phase:        phase,
			Metadata:     metadata,
		})
	}

	return
}

func (r *ThreatRepositorySQL) StatsByPhase(ctx context.Context, start, end time.Time) (output model.PieChartData, err error) {
	query := `SELECT
    		phase, COUNT(*)
			FROM threats
			WHERE seen_at >= $1 AND seen_at <= $2
			GROUP BY phase
			ORDER BY phase`

	rows, err := r.connection.Query(ctx, query, start, end)
	if err != nil {
		return
	}

	var dataset model.Dataset

	defer rows.Close()
	for rows.Next() {
		var phase string
		var total int

		err = rows.Scan(&phase, &total)
		if err != nil {
			return
		}

		output.Labels = append(output.Labels, phase)
		dataset.Data = append(dataset.Data, total)
	}
	dataset.Label = "Stats Threats By Phase"
	output.Datasets = append(output.Datasets, dataset)
	return output, err
}

func (r *ThreatRepositorySQL) StatsBySeverity(ctx context.Context, start, end time.Time) (output model.PieChartData, err error) {
	query := `SELECT
    		severity, COUNT(*)
			FROM stats_threat_by_severity
			WHERE _bucket >= $1 AND _bucket <= $2
			GROUP BY severity
			ORDER BY severity`

	rows, err := r.connection.Query(ctx, query, start, end)
	if err != nil {
		return
	}

	var dataset model.Dataset

	defer rows.Close()
	for rows.Next() {
		var severity string
		var total int

		err = rows.Scan(&severity, &total)
		if err != nil {
			return
		}

		output.Labels = append(output.Labels, severity)
		dataset.Data = append(dataset.Data, total)
	}
	dataset.Label = "Stats Threats By Severity"
	output.Datasets = append(output.Datasets, dataset)
	return output, err
}

func (r *ThreatRepositorySQL) TopHostAffected(ctx context.Context, start, end time.Time) (output model.BarChartData, err error) {
	query := `SELECT affected_host, COUNT(*) as c FROM threats
				WHERE seen_at >= $1 AND seen_at <= $2
				GROUP BY affected_host
				ORDER BY c
 				LIMIT 10`

	rows, err := r.connection.Query(ctx, query, start, end)
	if err != nil {
		return
	}

	var dataset model.Dataset
	defer rows.Close()
	for rows.Next() {
		var host string
		var count int

		err = rows.Scan(&host, &count)
		if err != nil {
			return
		}

		output.Labels = append(output.Labels, host)
		dataset.Data = append(dataset.Data, count)
	}
	dataset.Label = "Top 10 Affected Host"
	output.Datasets = append(output.Datasets, dataset)
	return output, err
}

func (r *ThreatRepositorySQL) TopAttacker(ctx context.Context, start, end time.Time) (output model.BarChartData, err error) {
	query := `SELECT attacker_host, COUNT(*) as c FROM threats
				WHERE seen_at >= $1 AND seen_at <= $2
				GROUP BY attacker_host
				ORDER BY c
 				LIMIT 10`

	rows, err := r.connection.Query(ctx, query, start, end)
	if err != nil {
		return
	}

	var dataset model.Dataset

	defer rows.Close()
	for rows.Next() {
		var host string
		var count int

		err = rows.Scan(&host, &count)
		if err != nil {
			return
		}

		output.Labels = append(output.Labels, host)
		dataset.Data = append(dataset.Data, count)
	}
	dataset.Label = "Top 10 Suspected Host"
	output.Datasets = append(output.Datasets, dataset)
	return output, err
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
		query := "INSERT INTO threats (id, created_at, seen_at, affected_host, attacker_host, conn_id, severity, confidence, phase, metadata)" +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
		batch.Queue(query,
			threat.ID,
			time.Now(),
			threat.SeenAt,
			threat.AffectedHost,
			threat.AttackerHost,
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
