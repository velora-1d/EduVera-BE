package model

import (
	"time"
)

type Tabungan struct {
	ID        string    `json:"id" goqu:"id" db:"id"`
	TenantID  string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	SantriID  string    `json:"santri_id" goqu:"santri_id" db:"santri_id"`
	Saldo     int64     `json:"saldo" goqu:"saldo" db:"saldo"`
	Status    string    `json:"status" goqu:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`

	// Joins
	NamaSantri string `json:"nama_santri" goqu:"skip" db:"nama_santri"`
	NIS        string `json:"nis" goqu:"skip" db:"nis"`
}

type TabunganMutasi struct {
	ID         string    `json:"id" goqu:"id" db:"id"`
	TabunganID string    `json:"tabungan_id" goqu:"tabungan_id" db:"tabungan_id"`
	TenantID   string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	Tipe       string    `json:"tipe" goqu:"tipe" db:"tipe"` // Debit, Kredit
	Nominal    int64     `json:"nominal" goqu:"nominal" db:"nominal"`
	Keterangan string    `json:"keterangan" goqu:"keterangan" db:"keterangan"`
	Petugas    string    `json:"petugas" goqu:"petugas" db:"petugas"`
	CreatedAt  time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
}
