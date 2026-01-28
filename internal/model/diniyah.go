package model

import (
	"time"
)

type DiniyahKitab struct {
	ID          string    `json:"id" goqu:"id" db:"id"`
	TenantID    string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	NamaKitab   string    `json:"nama_kitab" goqu:"nama_kitab" db:"nama_kitab"`
	BidangStudi string    `json:"bidang_studi" goqu:"bidang_studi" db:"bidang_studi"`
	Pengarang   string    `json:"pengarang" goqu:"pengarang" db:"pengarang"`
	Keterangan  string    `json:"keterangan" goqu:"keterangan" db:"keterangan"`
	CreatedAt   time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`
}
