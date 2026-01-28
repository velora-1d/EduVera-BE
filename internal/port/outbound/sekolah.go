package outbound_port

import "eduvera/internal/model"

type SekolahPort interface {
	GetSiswaByTenant(tenantID string) ([]model.Siswa, error)
	CreateSiswa(siswa model.Siswa) error
	GetGuruByTenant(tenantID string) ([]model.Guru, error)
	CreateGuru(guru model.Guru) error
	GetMapelByTenant(tenantID string) ([]model.Mapel, error)
}
