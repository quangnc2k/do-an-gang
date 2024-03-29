package model

import (
	"sync"
	"time"
)

type Alert struct {
	ID         string                 `json:"id"`
	CreatedAt  time.Time              `json:"created_at"`
	Details    map[string]interface{} `json:"details"`
	Resolved   bool                   `json:"resolved"`
	ResolvedAt *time.Time             `json:"resolved_at,omitempty"`
	ResolvedBy *string                `json:"resolved_by,omitempty"`
}

type AlertConfig struct {
	ID                string        `json:"id,omitempty"`
	Name              string        `json:"name,omitempty"`
	CreatedAt         time.Time     `json:"createdAt"`
	SeverityString    string        `json:"severity"`
	Severity          int           `json:"-"`
	Confidence        float64       `json:"confidence"`
	Recipients        []string      `json:"recipients,omitempty"`
	SuppressFor       time.Duration `json:"-"`
	SuppressForString string        `json:"suppress_for,omitempty"`

	mu        *sync.Mutex
	LastAlert time.Time `json:"-"`
}

func (c *AlertConfig) SetLock() {
	c.mu = new(sync.Mutex)
}

func (c *AlertConfig) Lock() {
	c.mu.Lock()
}

func (c *AlertConfig) Unlock() {
	c.mu.Unlock()
}
