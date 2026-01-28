package model

import (
	"time"
)

type PelanggaranAturan struct {
	ID        string    `json:"id" goqu:"skipinsert"`
	TenantID  string    `json:"tenant_id"`
	Judul     string    `json:"judul"`
	Kategori  string    `json:"kategori"`
	Poin      int       `json:"poin"`
	Level     string    `json:"level"` // Ringan, Sedang, Berat
	CreatedAt time.Time `json:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `json:"updated_at" goqu:"skipinsert"`
}

type PelanggaranSiswa struct {
	ID          string    `json:"id" goqu:"skipinsert"`
	TenantID    string    `json:"tenant_id"`
	SantriID    string    `json:"santri_id"`
	SantriNama  string    `json:"santri_nama" goqu:"skipinsert"` // Joined
	AturanID    *string   `json:"aturan_id"`
	AturanJudul string    `json:"aturan_judul" goqu:"skipinsert"` // Joined
	Tanggal     time.Time `json:"tanggal"`
	Poin        int       `json:"poin"`
	Keterangan  string    `json:"keterangan"`
	Status      string    `json:"status"` // Pending, Diproses, Selesai
	Sanksi      string    `json:"sanksi"`
	CreatedAt   time.Time `json:"created_at" goqu:"skipinsert"`
	UpdatedAt   time.Time `json:"updated_at" goqu:"skipinsert"`
}

type Perizinan struct {
	ID            string    `json:"id" goqu:"skipinsert"`
	TenantID      string    `json:"tenant_id"`
	SantriID      string    `json:"santri_id"`
	SantriNama    string    `json:"santri_nama" goqu:"skipinsert"` // Joined
	Tipe          string    `json:"tipe"`                          // Izin Pulang, Izin Keluar, Izin Sakit
	Alasan        string    `json:"alasan"`
	Dari          time.Time `json:"dari"`
	Sampai        time.Time `json:"sampai"`
	Status        string    `json:"status"` // Pending, Disetujui, Ditolak
	PenyetujuID   *string   `json:"penyetuju_id"`
	PenyetujuNama string    `json:"penyetuju_nama" goqu:"skipinsert"` // Joined
	CreatedAt     time.Time `json:"created_at" goqu:"skipinsert"`
	UpdatedAt     time.Time `json:"updated_at" goqu:"skipinsert"`
}
