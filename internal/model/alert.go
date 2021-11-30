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
