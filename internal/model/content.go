package model

import "time"

type Content struct {
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"` // Stored as string/json
	Type      string    `json:"type" db:"type"`   // 'text', 'json', 'image'
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ContentInput struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
	Type  string `json:"type" validate:"required"`
}
