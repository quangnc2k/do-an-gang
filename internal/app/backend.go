package app

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/app/backend"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"log"
)

func ServeBackend(ctx context.Context, addr string) (err error) {
	config.Init()

	persistance.LoadRepoContainerWithPgx(ctx)

	err = backend.ServeBackend(ctx, addr)
	if err != nil {
		log.Fatalln(err)
	}

	return
}
