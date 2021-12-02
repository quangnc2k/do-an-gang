package config

import (
	"log"

	"github.com/quangnc2k/do-an-gang/pkg/env"
)

var Env Environment

type Environment struct {
	PostgresHost     string `envconfig:"postgres_host"`
	PostgresPort     string `envconfig:"postgres_port"`
	PostgresUsername string `envconfig:"postgres_username"`
	PostgresPassword string `envconfig:"postgres_password"`
	PostgresDb       string `envconfig:"postgres_db"`

	RedisURL string `envconfig:"redis_url"`

	VTTApiKey         string `envconfig:"vtt_api_key"`
	VTTMaxFilerPerMin int    `envconfig:"vtt_max_file_per_min"`

	XForceUsername string `envconfig:"xforce_username"`
	XForcePassword string `envconfig:"xforce_password"`

	MigrateFilePath string `envconfig:"migrate_file_path"`

	JWTSecret string `envconfig:"jwt_secret"`
}

func Init() {
	var environment Environment
	err := env.ReadEnvVars(&environment)
	if err != nil {
		log.Fatalln("Unable to read env variables: ", err)
	}

	Env = environment
}
