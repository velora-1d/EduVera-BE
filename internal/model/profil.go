package model

import (
	"time"
)

type Profil struct {
	ID             string    `json:"id" goqu:"id" db:"id"`
	TenantID       string    `json:"tenant_id" goqu:"tenant_id" db:"tenant_id"`
	JenisPesantren string    `json:"jenis_pesantren" goqu:"jenis_pesantren" db:"jenis_pesantren"`
	Deskripsi      string    `json:"deskripsi" goqu:"deskripsi" db:"deskripsi"`
	Website        string    `json:"website" goqu:"website" db:"website"`
	EmailKontak    string    `json:"email_kontak" goqu:"email_kontak" db:"email_kontak"`
	NoTelpKontak   string    `json:"no_telp_kontak" goqu:"no_telp_kontak" db:"no_telp_kontak"`
	LogoURL        string    `json:"logo_url" goqu:"logo_url" db:"logo_url"`
	Curriculum     string    `json:"curriculum" goqu:"curriculum" db:"curriculum"` // K13, MERDEKA, PESANTREN
	CreatedAt      time.Time `json:"created_at" goqu:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" goqu:"updated_at" db:"updated_at"`

	// Joined Fields (from tenants)
	NamaPesantren string `json:"nama_pesantren" goqu:"skip" db:"nama_pesantren"` // from tenants.name
	Alamat        string `json:"alamat" goqu:"skip" db:"alamat"`                 // from tenants.address
}

type ProfilUpdate struct {
	NamaPesantren  string `json:"nama_pesantren"`
	JenisPesantren string `json:"jenis_pesantren"`
	Alamat         string `json:"alamat"`
	Deskripsi      string `json:"deskripsi"`
	Curriculum     string `json:"curriculum"` // K13, MERDEKA, PESANTREN
}
