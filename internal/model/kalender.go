package model

import (
	"time"
)

type KalenderEvent struct {
	ID          string    `json:"id" goqu:"id" db:"id"`
	TenantID    string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	Title       string    `json:"title" goqu:"title" db:"title"`
	StartDate   string    `json:"start_date" goqu:"start_date" db:"start_date"` // YYYY-MM-DD
	EndDate     string    `json:"end_date" goqu:"end_date" db:"end_date"`       // YYYY-MM-DD
	Category    string    `json:"category" goqu:"category" db:"category"`
	Description string    `json:"description" goqu:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`
}
