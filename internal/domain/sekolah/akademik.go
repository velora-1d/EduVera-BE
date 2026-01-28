package sekolah

import (
	"context"
	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

// Domain Interface
type AkademikDomain interface {
	// Siswa
	GetSiswaList(ctx context.Context, tenantID string) ([]model.Siswa, error)
	CreateSiswa(ctx context.Context, siswa model.Siswa) error

	// Guru
	GetGuruList(ctx context.Context, tenantID string) ([]model.Guru, error)
	CreateGuru(ctx context.Context, guru model.Guru) error

	// Mapel
	GetMapelList(ctx context.Context, tenantID string) ([]model.Mapel, error)
}

// Implementation
type akademikDomain struct {
	databasePort outbound_port.DatabasePort
}

func NewAkademikDomain(databasePort outbound_port.DatabasePort) AkademikDomain {
	return &akademikDomain{
		databasePort: databasePort,
	}
}

// ------ Siswa Implementation ------

func (d *akademikDomain) GetSiswaList(ctx context.Context, tenantID string) ([]model.Siswa, error) {
	return d.databasePort.Sekolah().GetSiswaByTenant(tenantID)
}

func (d *akademikDomain) CreateSiswa(ctx context.Context, siswa model.Siswa) error {
	return d.databasePort.Sekolah().CreateSiswa(siswa)
}

// ------ Guru Implementation ------

func (d *akademikDomain) GetGuruList(ctx context.Context, tenantID string) ([]model.Guru, error) {
	return d.databasePort.Sekolah().GetGuruByTenant(tenantID)
}

func (d *akademikDomain) CreateGuru(ctx context.Context, guru model.Guru) error {
	return d.databasePort.Sekolah().CreateGuru(guru)
}

// ------ Mapel Implementation ------

func (d *akademikDomain) GetMapelList(ctx context.Context, tenantID string) ([]model.Mapel, error) {
	return d.databasePort.Sekolah().GetMapelByTenant(tenantID)
}
