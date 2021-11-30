package app

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/app/process"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
)

func Process(ctx context.Context) (err error) {
	config.Init()

	persistance.LoadRepoContainerWithPgx(ctx)
	persistance.LoadQueueWithRedis(ctx)

	persistance.LoadBlacklistEngine(ctx)
	persistance.LoadFileEngine(ctx)
	persistance.LoadIPEngine(ctx)

	return process.Run(ctx)
}
