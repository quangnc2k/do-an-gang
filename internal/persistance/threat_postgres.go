package persistance

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (r *ThreatRepositorySQL) RecentAffected(ctx context.Context) (hosts map[string]time.Time, err error) {
	query := `SELECT affected_host, seen_at FROM threats ORDER BY seen_at LIMIT 50`

	rows, err := r.connection.Query(ctx, query)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var host string
		var t time.Time

		err = rows.Scan(&host, &t)
		if err != nil {
			log.Println(err)
			continue
		}

		hosts[host] = t
	}

	return
}

func (r *ThreatRepositorySQL) RecentAttackByPhase(ctx context.Context) (phases map[string]time.Time, err error) {
	query := `SELECT phase, seen_at FROM threats ORDER BY seen_at LIMIT 50`

	rows, err := r.connection.Query(ctx, query)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var ph string
		var t time.Time

		err = rows.Scan(&ph, &t)
		if err != nil {
			log.Println(err)
			continue
		}

		phases[ph] = t
	}

	return
}

func (r *ThreatRepositorySQL) Overview(ctx context.Context) (total, recent, numOfHost int64) {
	query := `SELECT COUNT(*) FROM threats`
	err := r.connection.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return
	}

	query = `SELECT COUNT(*) FROM threats WHERE seen_at >= now() - INTERVAL '7 day' AND seen_at <= now()`
	err = r.connection.QueryRow(ctx, query).Scan(&recent)
	if err != nil {
		return
	}

	query = `SELECT COUNT(DISTINCT affected_host) FROM threats`
	err = r.connection.QueryRow(ctx, query).Scan(&numOfHost)
	if err != nil {
		return
	}

	return
}

func NewThreatSQLRepo(ctx context.Context, conn *pgxpool.Pool) (ThreatRepository, error) {
	if conn == nil {
		return nil, errors.New("invalid sql connection")
	}
	return &ThreatRepositorySQL{connection: conn}, nil
}

func (r *ThreatRepositorySQL) HistogramAffected(ctx context.Context, width string, start, end time.Time) (output model.LineChartData, err error) {
	query1 := `SELECT affected_host FROM stats_threat_by_affected
						WHERE _bucket >= $1 AND _bucket <= $2
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
		query := `SELECT time_bucket_gapfill($1, _bucket) as t, SUM(count) as total
					FROM stats_threat_by_affected
					WHERE _bucket >= $2 AND _bucket <= $3 AND affected_host = $4
					GROUP BY t, affected_host
					ORDER BY t ASC`

		rows, err := r.connection.Query(ctx, query, width, start, end, host)
		if err != nil {
			return output, err
		}

		var dataset model.Dataset

		for rows.Next() {
			var t time.Time
			var total *int
			var emptyTotal int

			err := rows.Scan(&t, &total)
			if err != nil {
				rows.Close()
				return output, err
			}

			if i == 0 {
				output.Labels = append(output.Labels, t.String())
			}

			if total == nil {
				total = &emptyTotal
			}

			dataset.Data = append(dataset.Data, *total)
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
			sq.Like{"severity": search},
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
		var id, affectedHost, attackHost, severity, phase string
		var confidence float64
		var metadata map[string]interface{}

		err = rows.Scan(&id, &createdAt, &seenAt, &affectedHost, &attackHost, &confidence, &severity, &phase, &metadata)
		if err != nil {
			log.Println(err)
			continue
		}

		output.Data = append(output.Data, model.Threat{
			ID:             id,
			CreatedAt:      createdAt,
			SeenAt:         seenAt,
			AffectedHost:   affectedHost,
			AttackerHost:   attackHost,
			Confidence:     confidence,
			SeverityString: severity,
			Phase:          phase,
			Metadata:       metadata,
		})
	}

	return
}

func (r *ThreatRepositorySQL) StatsByPhase(ctx context.Context, start, end time.Time) (output model.PieChartData, err error) {
	query := `SELECT
    		phase, SUM(count)
			FROM stats_threat_by_phase
			WHERE _bucket >= $1 AND _bucket <= $2
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
    		severity, SUM(count)
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

func (r *ThreatRepositorySQL) StoreThreatInBatch(ctx context.Context, threats []model.Threat) (err error) {
	batch := &pgx.Batch{}

	if len(threats) == 0 {
		return
	}

	for _, threat := range threats {
		var severity string
		if threat.Severity > 8 {
			severity = "critical"
		} else if threat.Severity > 5 {
			severity = "high"
		} else if threat.Severity > 2 {
			severity = "medium"
		} else {
			severity = "low"
		}
		query := "INSERT INTO threats (id, created_at, seen_at, affected_host, attacker_host, conn_id, severity, confidence, phase, metadata)" +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
		batch.Queue(query,
			uuid.New().String(),
			time.Now(),
			threat.SeenAt,
			threat.AffectedHost,
			threat.AttackerHost,
			threat.ConnID,
			severity,
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
