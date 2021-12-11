package services

import (
	"context"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"gopkg.in/gomail.v2"
	"time"
)

var alertStore AlertStore

type AlertStore struct {
	alertConfigStore []model.AlertConfig
}

func (s *AlertStore) InitAlertStore(ctx context.Context) (err error) {
	configs, err := persistance.GetRepoContainer().AlertConfigRepository.GetAll(ctx)
	if err != nil {
		return
	}

	alertStore.alertConfigStore = configs
	return
}

func (s *AlertStore) CheckAlert(ctx context.Context, threat model.Threat) (err error) {
	for i, alertConfig := range alertStore.alertConfigStore {
		if threat.Severity >= float64(alertConfig.Severity) &&
			alertConfig.LastAlert.Add(alertConfig.SuppressFor).Before(time.Now()) {
			err = s.SendAlertEmail(ctx, &alertConfig)
			if err != nil {
				return
			}

			alertStore.alertConfigStore[i].Lock()
			alertStore.alertConfigStore[i].LastAlert = time.Now()
			alertStore.alertConfigStore[i].Unlock()
		}
	}

	return err
}

func (s *AlertStore) SendAlertEmail(ctx context.Context, c *model.AlertConfig) (err error) {
	from := config.Env.MailUser
	pass := config.Env.MailPassword
	to := c.Recipients

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "Alert", c.Name, "was triggered")
	m.SetBody("text/html", "Needs something here!")

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)

	// Send the email to Bob, Cora and Dan.
	if err = d.DialAndSend(m); err != nil {
		return
	}

	return
}
