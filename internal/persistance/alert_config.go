package persistance

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"sync"
)

type AlertConfigRepository interface {
	GetAll(ctx context.Context) (configs []model.AlertConfig, err error)
	Create(ctx context.Context, config *model.AlertConfig) (err error)
	FindOneByID(ctx context.Context, id string) (config model.AlertConfig, err error)
	UpdateOneByID(ctx context.Context, config model.AlertConfig, id string) (err error)
	DeleteOneByID(ctx context.Context, id string) (err error)

	ListenAndUpdate(ctx context.Context, configs []model.AlertConfig, l *sync.RWMutex) (err error)
}
