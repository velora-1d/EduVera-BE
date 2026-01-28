package model

import (
	"time"
)

type TahfidzSetoran struct {
	ID        string    `json:"id" goqu:"id" db:"id"`
	TenantID  string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	SantriID  string    `json:"santri_id" goqu:"santri_id" db:"santri_id"`
	UstadzID  *string   `json:"ustadz_id" goqu:"ustadz_id" db:"ustadz_id"`
	Tanggal   time.Time `json:"tanggal" goqu:"tanggal" db:"tanggal"`
	Juz       int       `json:"juz" goqu:"juz" db:"juz"`
	Surah     string    `json:"surah" goqu:"surah" db:"surah"`
	AyatAwal  int       `json:"ayat_awal" goqu:"ayat_awal" db:"ayat_awal"`
	AyatAkhir int       `json:"ayat_akhir" goqu:"ayat_akhir" db:"ayat_akhir"`
	Tipe      string    `json:"tipe" goqu:"tipe" db:"tipe"`             // Ziyadah, Murajaah
	Kualitas  string    `json:"kualitas" goqu:"kualitas" db:"kualitas"` // Lancar, Kurang, Ulang
	Catatan   string    `json:"catatan" goqu:"catatan" db:"catatan"`
	CreatedAt time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`

	// Joins
	SantriNama string `json:"santri_nama" goqu:"skipinsert" db:"santri_nama"`
	UstadzNama string `json:"ustadz_nama" goqu:"skipinsert" db:"ustadz_nama"`
}
