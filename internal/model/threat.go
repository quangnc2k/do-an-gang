package model

import "time"

type Threat struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	SeenAt          time.Time `json:"seen_at"`
	SourceHost      string    `json:"src_host"`
	DestinationHost string    `json:"dst_host,omitempty"`
	ConnID          string    `json:"conn_id,omitempty"`
	Confidence      float64   `json:"confidence"`
	Severity        int       `json:"severity"`
	Phase           string    `json:"phase"`

	Metadata interface{} `json:"metadata,omitempty"`
}
