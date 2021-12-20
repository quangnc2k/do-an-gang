package persistance

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"log"
	"sync"
)

type AlertConfigRepositorySQL struct {
	connection *pgxpool.Pool
}

func (r *AlertConfigRepositorySQL) ListenAndUpdate(ctx context.Context, configs []model.AlertConfig, l *sync.RWMutex) (err error) {
	c, err := r.connection.Acquire(ctx)
	if err != nil {
		return
	}

	_, err = c.Exec(ctx, "LISTEN changes;")
	if err != nil {
		c.Release()
		return
	}

	go func() {
		defer c.Release()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				var payload struct {
					Table  string      `json:"table"`
					Action string      `json:"action"`
					Data   interface{} `json:"data"`
				}

				notification, err := c.Conn().WaitForNotification(ctx)
				if err == context.Canceled {
					return
				}

				if err != nil {
					return
				}

				b := []byte(notification.Payload)

				err = json.Unmarshal(b, &payload)
				if err != nil {
					//TODO add log
					log.Println("Listening alert_configs", err)
					continue
				}

				if payload.Table != "alert_configs" {
					continue
				}

				var listenedConfig model.AlertConfig

				valueByte, err := json.Marshal(payload.Data)
				if err != nil {
					//TODO add log
					log.Println("Listening alert_configs", err)
					continue
				}

				err = json.Unmarshal(valueByte, &listenedConfig)
				if err != nil {
					//TODO add log
					log.Println("Listening alert_configs", err)
					continue
				}

				l.Lock()

				if payload.Action == "DELETE" {
					for i, config := range configs {
						if config.ID == listenedConfig.ID {
							configs[i] = configs[len(configs)-1]
							configs = configs[:len(configs)-1]
						}
					}
				} else if payload.Action == "UPDATE" {
					for i, config := range configs {
						if config.ID == listenedConfig.ID {
							listenedConfig.SetLock()
							configs[i] = listenedConfig
						}
					}
				} else if payload.Action == "INSERT" {
					listenedConfig.SetLock()
					configs = append(configs, listenedConfig)
				}

				l.Unlock()
			}
		}
	}()
	return
}

func (r *AlertConfigRepositorySQL) GetAll(ctx context.Context) (configs []model.AlertConfig, err error) {
	query := `SELECT id, name, created_at, severity, confidence, recipients, suppress_for
				FROM alert_configs`

	rows, err := r.connection.Query(ctx, query)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var config model.AlertConfig

		err = rows.Scan(&config.ID, &config.Name, &config.CreatedAt, &config.SeverityString, &config.Confidence, &config.Recipients, &config.SuppressFor)
		if err != nil {
			log.Println("get all configs", err)
			continue
		}

		if config.SeverityString == "low" {
			config.Severity = 0
		} else if config.SeverityString == "medium" {
			config.Severity = 4
		} else if config.SeverityString == "high" {
			config.Severity = 7
		} else if config.SeverityString == "critical" {
			config.Severity = 9
		}

		config.SetLock()

		configs = append(configs, config)
	}

	return
}

func (r *AlertConfigRepositorySQL) Create(ctx context.Context, config *model.AlertConfig) (err error) {
	query := `INSERT INTO alert_configs (id, name, created_at, severity, confidence, recipients, suppress_for)
				VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err = r.connection.Exec(ctx, query,
		config.ID,
		config.Name,
		config.CreatedAt,
		config.SeverityString,
		config.Confidence,
		config.Recipients,
		config.SuppressFor,
	)
	if err != nil {
		return
	}

	return
}

func (r *AlertConfigRepositorySQL) FindOneByID(ctx context.Context, id string) (config model.AlertConfig, err error) {
	query := `SELECT id, name, created_at, severity, confidence, recipients, suppress_for
				FROM alert_configs
				WHERE id = $1`

	err = r.connection.QueryRow(ctx, query, id).
		Scan(&config.ID, &config.Name, &config.CreatedAt, &config.SeverityString, &config.Confidence, &config.Recipients, &config.SuppressFor)
	if err != nil {
		return
	}

	return
}

func (r *AlertConfigRepositorySQL) UpdateOneByID(ctx context.Context, config model.AlertConfig, id string) (err error) {
	_, err = r.FindOneByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.New("invalid id")
		}
		return
	}

	query := `UPDATE alert_configs
				SET name  = $1,
					severity = $2,
					confidence = $3,
					recipients = $4,
					suppress_for = $5
				WHERE id = $6`
	_, err = r.connection.Exec(ctx, query,
		config.Name,
		config.SeverityString,
		config.Confidence,
		config.Recipients,
		config.SuppressFor,
		id,
	)
	if err != nil {
		return
	}

	return
}

func (r *AlertConfigRepositorySQL) DeleteOneByID(ctx context.Context, id string) (err error) {
	_, err = r.FindOneByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.New("invalid id")
		}
		return
	}

	query := `DELETE FROM alert_configs
				WHERE id = $1`
	_, err = r.connection.Exec(ctx, query, id)
	if err != nil {
		return
	}

	return
}

func NewAlertConfigSQLRepo(ctx context.Context, conn *pgxpool.Pool) (AlertConfigRepository, error) {
	if conn == nil {
		return nil, errors.New("invalid sql connection")
	}
	return &AlertConfigRepositorySQL{connection: conn}, nil
}
