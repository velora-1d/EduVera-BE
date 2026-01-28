package model

import "time"

type Asrama struct {
	ID        string    `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Nama      string    `json:"nama" db:"nama"`
	Jenis     string    `json:"jenis" db:"jenis"` // Putra, Putri
	MusyrifID *string   `json:"musyrif_id" db:"musyrif_id"`
	Musyrif   string    `json:"musyrif_nama" db:"musyrif_nama"` // Populated via join
	Status    string    `json:"status" db:"status"`
	Kapasitas int       `json:"kapasitas" db:"kapasitas"` // Calculated
	Terisi    int       `json:"terisi" db:"terisi"`       // Calculated
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Kamar struct {
	ID         string    `json:"id" db:"id"`
	TenantID   string    `json:"tenant_id" db:"tenant_id"`
	AsramaID   string    `json:"asrama_id" db:"asrama_id"`
	AsramaNama string    `json:"asrama_nama" db:"asrama_nama"` // Populated via join
	Nomor      string    `json:"nomor" db:"nomor"`
	Kapasitas  int       `json:"kapasitas" db:"kapasitas"`
	Terisi     int       `json:"terisi" db:"terisi"` // Calculated
	Status     string    `json:"status" db:"status"` // Tersedia, Penuh
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Penempatan struct {
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	SantriID     string    `json:"santri_id" db:"santri_id"`
	SantriNama   string    `json:"santri_nama" db:"santri_nama"` // Populated via join
	KamarID      string    `json:"kamar_id" db:"kamar_id"`
	KamarNomor   string    `json:"kamar_nomor" db:"kamar_nomor"` // Populated via join
	AsramaNama   string    `json:"asrama_nama" db:"asrama_nama"` // Populated via join
	TanggalMasuk time.Time `json:"tanggal_masuk" db:"tanggal_masuk"`
	Status       string    `json:"status" db:"status"`
	Keterangan   string    `json:"keterangan" db:"keterangan"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
