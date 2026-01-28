package model

import (
	"time"
)

type RaporPeriode struct {
	ID           string    `json:"id" goqu:"id" db:"id"`
	TenantID     string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	Nama         string    `json:"nama" goqu:"nama" db:"nama"`
	TanggalMulai string    `json:"tanggal_mulai" goqu:"tanggal_mulai" db:"tanggal_mulai"`
	TanggalAkhir string    `json:"tanggal_akhir" goqu:"tanggal_akhir" db:"tanggal_akhir"`
	IsActive     bool      `json:"is_active" goqu:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`
}

type Rapor struct {
	ID               string    `json:"id" goqu:"id" db:"id"`
	TenantID         string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	PeriodeID        string    `json:"periode_id" goqu:"periode_id" db:"periode_id"`
	SantriID         string    `json:"santri_id" goqu:"santri_id" db:"santri_id"`
	Status           string    `json:"status" goqu:"status" db:"status"`
	CatatanWaliKelas string    `json:"catatan_wali_kelas" goqu:"catatan_wali_kelas" db:"catatan_wali_kelas"`
	CreatedAt        time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`

	// Joins
	NamaSantri  string       `json:"nama_santri" goqu:"skip" db:"nama_santri"`
	NamaPeriode string       `json:"nama_periode" goqu:"skip" db:"nama_periode"`
	NilaiList   []RaporNilai `json:"nilai_list" goqu:"skip" db:"skip"`
}

type RaporNilai struct {
	ID         string    `json:"id" goqu:"id" db:"id"`
	RaporID    string    `json:"rapor_id" goqu:"rapor_id" db:"rapor_id"`
	Kategori   string    `json:"kategori" goqu:"kategori" db:"kategori"`
	Jenis      string    `json:"jenis" goqu:"jenis" db:"jenis"`
	Nilai      string    `json:"nilai" goqu:"nilai" db:"nilai"`
	Keterangan string    `json:"keterangan" goqu:"keterangan" db:"keterangan"`
	CreatedAt  time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`
}
