package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

func InitAlertStore(ctx context.Context) (err error) {
	configs, err := persistance.GetRepoContainer().AlertConfigRepository.GetAll(ctx)
	if err != nil {
		return
	}

	alertStore.alertConfigStore = configs
	return
}

func (s *AlertStore) CheckAlert(ctx context.Context, threat model.Threat) (err error) {
	for i, _ := range alertStore.alertConfigStore {
		alertStore.alertConfigStore[i].Lock()
		if threat.Severity >= float64(alertStore.alertConfigStore[i].Severity) &&
			alertStore.alertConfigStore[i].LastAlert.Add(alertStore.alertConfigStore[i].SuppressFor).Before(time.Now()) {
			var str = ""
			m := make(map[string]interface{})

			alertStore.alertConfigStore[i].LastAlert = time.Now()

			data, err := json.Marshal(threat)
			if err != nil {
				alertStore.alertConfigStore[i].Unlock()
				return err
			}

			err = json.Unmarshal(data, &m)
			if err != nil {
				alertStore.alertConfigStore[i].Unlock()
				return err
			}

			err = persistance.GetRepoContainer().AlertRepository.Create(ctx,
				model.Alert{
					ID:         uuid.NewString(),
					CreatedAt:  time.Now(),
					Details:    m,
					Resolved:   false,
					ResolvedAt: &time.Time{},
					ResolvedBy: &str,
				})
			if err != nil {
				alertStore.alertConfigStore[i].Unlock()
				return err
			}

			err = s.SendAlertEmail(ctx, &alertStore.alertConfigStore[i], threat)
			if err != nil {
				alertStore.alertConfigStore[i].Unlock()
				return err
			}
		}
		alertStore.alertConfigStore[i].Unlock()
	}

	return err
}

func (s *AlertStore) SendAlertEmail(ctx context.Context, c *model.AlertConfig, threat model.Threat) (err error) {
	from := config.Env.MailUser
	pass := config.Env.MailPassword
	to := c.Recipients

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "Alert", c.Name, "was triggered")
	m.SetBody("text/html", fmt.Sprintf("%+v", threat))

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)

	// Send the email to Bob, Cora and Dan.
	if err = d.DialAndSend(m); err != nil {
		return
	}

	return
}
