package model

import (
	"time"
)

// Student types
const (
	StudentTypeSiswa  = "siswa"
	StudentTypeSantri = "santri"
	StudentTypeBoth   = "both"
)

// Student status
const (
	StudentStatusActive      = "active"
	StudentStatusGraduated   = "graduated"
	StudentStatusDroppedOut  = "dropped_out"
	StudentStatusTransferred = "transferred"
)

// Student jenjang (education levels)
var ValidJenjangs = []string{
	"TK", "SD", "MI", "SMP", "MTs", "SMA", "MA", "SMK",
}

type Student struct {
	ID       string `json:"id" db:"id"`
	TenantID string `json:"tenant_id" db:"tenant_id"`

	// Basic Info
	Name       string     `json:"name" db:"name"`
	Nickname   *string    `json:"nickname,omitempty" db:"nickname"`
	Gender     *string    `json:"gender,omitempty" db:"gender"`
	BirthPlace *string    `json:"birth_place,omitempty" db:"birth_place"`
	BirthDate  *time.Time `json:"birth_date,omitempty" db:"birth_date"`
	PhotoURL   *string    `json:"photo_url,omitempty" db:"photo_url"`

	// NIS (multiple types supported)
	NIS          *string `json:"nis,omitempty" db:"nis"`
	NISN         *string `json:"nisn,omitempty" db:"nisn"`
	NISPesantren *string `json:"nis_pesantren,omitempty" db:"nis_pesantren"`

	// Contact & Address
	Address *string `json:"address,omitempty" db:"address"`
	Phone   *string `json:"phone,omitempty" db:"phone"`

	// Parent/Guardian Info
	FatherName       *string `json:"father_name,omitempty" db:"father_name"`
	FatherPhone      *string `json:"father_phone,omitempty" db:"father_phone"`
	FatherOccupation *string `json:"father_occupation,omitempty" db:"father_occupation"`
	MotherName       *string `json:"mother_name,omitempty" db:"mother_name"`
	MotherPhone      *string `json:"mother_phone,omitempty" db:"mother_phone"`
	MotherOccupation *string `json:"mother_occupation,omitempty" db:"mother_occupation"`
	GuardianName     *string `json:"guardian_name,omitempty" db:"guardian_name"`
	GuardianPhone    *string `json:"guardian_phone,omitempty" db:"guardian_phone"`
	GuardianRelation *string `json:"guardian_relation,omitempty" db:"guardian_relation"`

	// Type & Classification
	Type    string  `json:"type" db:"type"`
	Jenjang *string `json:"jenjang,omitempty" db:"jenjang"`

	// Academic Relations
	KelasID *string `json:"kelas_id,omitempty" db:"kelas_id"`

	// Pesantren Relations
	KamarID    *string `json:"kamar_id,omitempty" db:"kamar_id"`
	IsMukim    bool    `json:"is_mukim" db:"is_mukim"`
	TahunMasuk *int    `json:"tahun_masuk,omitempty" db:"tahun_masuk"`

	// Status
	Status     string     `json:"status" db:"status"`
	EntryDate  *time.Time `json:"entry_date,omitempty" db:"entry_date"`
	ExitDate   *time.Time `json:"exit_date,omitempty" db:"exit_date"`
	ExitReason *string    `json:"exit_reason,omitempty" db:"exit_reason"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type StudentInput struct {
	Name       string `json:"name" validate:"required"`
	Nickname   string `json:"nickname,omitempty"`
	Gender     string `json:"gender,omitempty"`
	BirthPlace string `json:"birth_place,omitempty"`
	BirthDate  string `json:"birth_date,omitempty"` // Format: YYYY-MM-DD

	NIS          string `json:"nis,omitempty"`
	NISN         string `json:"nisn,omitempty"`
	NISPesantren string `json:"nis_pesantren,omitempty"`

	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`

	FatherName       string `json:"father_name,omitempty"`
	FatherPhone      string `json:"father_phone,omitempty"`
	FatherOccupation string `json:"father_occupation,omitempty"`
	MotherName       string `json:"mother_name,omitempty"`
	MotherPhone      string `json:"mother_phone,omitempty"`
	MotherOccupation string `json:"mother_occupation,omitempty"`
	GuardianName     string `json:"guardian_name,omitempty"`
	GuardianPhone    string `json:"guardian_phone,omitempty"`
	GuardianRelation string `json:"guardian_relation,omitempty"`

	Type    string `json:"type" validate:"required,oneof=siswa santri both"`
	Jenjang string `json:"jenjang,omitempty"`

	KelasID    string `json:"kelas_id,omitempty"`
	KamarID    string `json:"kamar_id,omitempty"`
	IsMukim    bool   `json:"is_mukim,omitempty"`
	TahunMasuk int    `json:"tahun_masuk,omitempty"`

	Status string `json:"status,omitempty"`
}

type StudentFilter struct {
	TenantID string   `json:"tenant_id"`
	IDs      []string `json:"ids,omitempty"`
	Type     string   `json:"type,omitempty"` // siswa, santri, both
	Jenjang  string   `json:"jenjang,omitempty"`
	KelasID  string   `json:"kelas_id,omitempty"`
	KamarID  string   `json:"kamar_id,omitempty"`
	IsMukim  *bool    `json:"is_mukim,omitempty"`
	Status   string   `json:"status,omitempty"`
	Search   string   `json:"search,omitempty"` // Search by name, NIS, NISN
}
