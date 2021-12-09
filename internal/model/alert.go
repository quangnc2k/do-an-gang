package model

import "time"

type Alert struct {
	ID         string                 `json:"id"`
	CreatedAt  time.Time              `json:"created_at"`
	Details    map[string]interface{} `json:"details"`
	Resolved   bool                   `json:"resolved"`
	ResolvedAt *time.Time             `json:"resolved_at,omitempty"`
	ResolvedBy *string                `json:"resolved_by,omitempty"`
}

type AlertConfig struct {
	ID             string        `json:"id,omitempty"`
	Name           string        `json:"name,omitempty"`
	CreatedAt      time.Time     `json:"createdAt"`
	SeverityString string        `json:"severity"`
	Severity       int           `json:"-"`
	Confidence     float64       `json:"confidence,omitempty"`
	Recipients     []string      `json:"recipient,omitempty"`
	SuppressFor    time.Duration `json:"suppressFor,omitempty"`
}
