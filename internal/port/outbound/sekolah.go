package outbound_port

import "prabogo/internal/model"

type SekolahPort interface {
	GetSiswaByTenant(tenantID string) ([]model.Siswa, error)
	CreateSiswa(siswa model.Siswa) error
	GetGuruByTenant(tenantID string) ([]model.Guru, error)
	CreateGuru(guru model.Guru) error
	GetMapelByTenant(tenantID string) ([]model.Mapel, error)
	GetKelasByTenant(tenantID string) ([]model.Kelas, error)
	CreateKelas(kelas *model.Kelas) error

	// Asrama
	GetAsramaByTenant(tenantID string) ([]model.Asrama, error)
	CreateAsrama(asrama *model.Asrama) error
	GetKamarByAsrama(tenantID, asramaID string) ([]model.Kamar, error)
	CreateKamar(kamar *model.Kamar) error
	GetPenempatanByTenant(tenantID string) ([]model.Penempatan, error)
	CreatePenempatan(penempatan *model.Penempatan) error

	// Kepesantrenan
	GetPelanggaranAturan(tenantID string) ([]model.PelanggaranAturan, error)
	CreatePelanggaranAturan(m *model.PelanggaranAturan) error
	GetPelanggaranSiswa(tenantID string) ([]model.PelanggaranSiswa, error)
	CreatePelanggaranSiswa(m *model.PelanggaranSiswa) error
	GetPerizinan(tenantID string) ([]model.Perizinan, error)
	CreatePerizinan(m *model.Perizinan) error

	// Tahfidz
	GetTahfidzSetoran(tenantID string) ([]model.TahfidzSetoran, error)
	CreateTahfidzSetoran(m *model.TahfidzSetoran) error

	// Diniyah
	GetDiniyahKitab(tenantID string) ([]model.DiniyahKitab, error)
	CreateDiniyahKitab(m *model.DiniyahKitab) error

	// Rapor
	GetRaporList(tenantID string) ([]model.Rapor, error)
	CreateRapor(m *model.Rapor) error

	// Tabungan
	GetTabunganList(tenantID string) ([]model.Tabungan, error)
	CreateTabunganMutasi(m *model.TabunganMutasi) error

	// Kalender
	GetKalenderEvents(tenantID string) ([]model.KalenderEvent, error)
	CreateKalenderEvent(m *model.KalenderEvent) error

	// Profil
	GetProfil(tenantID string) (*model.Profil, error)
	UpdateProfil(tenantID string, m *model.ProfilUpdate) error

	// Laporan
	GetReportData(tenantID string, req model.ReportRequest) ([]model.ReportData, error)
}
