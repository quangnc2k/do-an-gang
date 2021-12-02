package persistance

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/quangnc2k/do-an-gang/internal/config"
)

var repoContainer RepositoryContainer

type RepositoryContainer struct {
	ThreatRepository ThreatRepository
	UserRepository   UserRepository
	AlertRepository  AlertRepository
}

func GetRepoContainer() RepositoryContainer {
	return repoContainer
}

func LoadRepoContainerWithPgx(ctx context.Context) {
	e := config.Env
	cfg, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		e.PostgresUsername, e.PostgresPassword, e.PostgresHost, e.PostgresPort, e.PostgresDb, "disable",
	))

	if err != nil {
		log.Fatalln("Cannot init repositories", err)
	}

	cfg.MaxConns = 10

	conn, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatalln("Cannot init repositories", err)
	}

	threatRepo, err := NewThreatSQLRepo(ctx, conn)
	if err != nil {
		log.Fatalln("Cannot init threat repositories", err)
	}

	userRepo, err := NewUserSQLRepo(ctx, conn)
	if err != nil {
		log.Fatalln("Cannot init user repositories", err)
	}

	alertRepo, err := NewAlertSQLRepo(ctx, conn)
	if err != nil {
		log.Fatalln("Cannot init alert repositories", err)
	}

	//Load database global var with this
	repoContainer = RepositoryContainer{
		ThreatRepository: threatRepo,
		UserRepository:   userRepo,
		AlertRepository:  alertRepo,
	}
}
