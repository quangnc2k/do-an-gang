package persistance

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/persistance/filters"
)

var BlacklistEngine Filter

var FileEngine Filter

var IPEngine Filter

type Filter interface {
	Check(ctx context.Context, resource string) (marked bool, credit float64, extraResource interface{}, err error)
}

func LoadBlacklistEngine(ctx context.Context) {
	BlacklistEngine = filters.InitFedoroEngine(ctx)
}

func LoadFileEngine(ctx context.Context) {
	FileEngine = filters.InitVTTEngine(ctx)
}

func LoadIPEngine(ctx context.Context) {
	IPEngine = filters.InitXForceEngine(ctx)
}
