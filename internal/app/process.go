package app

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/app/backend"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"log"
)

func Process(ctx context.Context) (err error) {
	config.Init()

	persistance.LoadRepoContainerWithPgx(ctx)
	//persistance.LoadQueueWithRedis(ctx)
	//
	//persistance.LoadBlacklistEngine(ctx)
	//persistance.LoadFileEngine(ctx)
	//persistance.LoadIPEngine(ctx)

	err = backend.ServeBackend(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return
	//return process.Run(ctx)
}
