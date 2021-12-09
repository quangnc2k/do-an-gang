package services

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
)

var alertConfigStore []model.AlertConfig

func InitAlertStore(ctx context.Context) (err error) {
	configs, err := persistance.GetRepoContainer().AlertConfigRepository.GetAll(ctx)
	if err != nil {
		return
	}

	alertConfigStore = configs
	return
}

func CheckAlert(ctx context.Context, threat model.Threat) (matched bool, err error) {
	for _, config := range alertConfigStore {
		if threat.Severity >= config.Severity {
		}
	}

	return false, err
}
