package env

import (
	"errors"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const EnvPrefix = "atd"

func ReadEnvVars(args ...interface{}) error {
	if len(args) == 1 {
		return envconfig.Process(EnvPrefix, args[0])
	}

	return errors.New("invalid env params")
}

func LoadEnv(file string) error {
	if file == "" {
		return godotenv.Load()
	}

	return godotenv.Load(file)
}
