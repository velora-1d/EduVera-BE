package sekolah

import (
	"context"
	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
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

	// Kelas (Marhalah)
	GetKelasList(ctx context.Context, tenantID string) ([]model.Kelas, error)
	CreateKelas(ctx context.Context, tenantID string, kelas *model.Kelas) error

	// Asrama
	GetAsramaList(ctx context.Context, tenantID string) ([]model.Asrama, error)
	CreateAsrama(ctx context.Context, tenantID string, asrama *model.Asrama) error
	GetKamarList(ctx context.Context, tenantID, asramaID string) ([]model.Kamar, error)
	CreateKamar(ctx context.Context, tenantID string, kamar *model.Kamar) error
	GetPenempatanList(ctx context.Context, tenantID string) ([]model.Penempatan, error)
	CreatePenempatan(ctx context.Context, tenantID string, penempatan *model.Penempatan) error

	// Kepesantrenan
	GetPelanggaranAturanList(ctx context.Context, tenantID string) ([]model.PelanggaranAturan, error)
	CreatePelanggaranAturan(ctx context.Context, tenantID string, m *model.PelanggaranAturan) error
	GetPelanggaranSiswaList(ctx context.Context, tenantID string) ([]model.PelanggaranSiswa, error)
	CreatePelanggaranSiswa(ctx context.Context, tenantID string, m *model.PelanggaranSiswa) error
	GetPerizinanList(ctx context.Context, tenantID string) ([]model.Perizinan, error)
	CreatePerizinan(ctx context.Context, tenantID string, m *model.Perizinan) error
	// Tahfidz
	GetTahfidzSetoranList(ctx context.Context, tenantID string) ([]model.TahfidzSetoran, error)
	CreateTahfidzSetoran(ctx context.Context, tenantID string, m *model.TahfidzSetoran) error
	// Diniyah
	GetDiniyahKitabList(ctx context.Context, tenantID string) ([]model.DiniyahKitab, error)
	CreateDiniyahKitab(ctx context.Context, tenantID string, m *model.DiniyahKitab) error
	// Rapor
	GetRaporList(ctx context.Context, tenantID string) ([]model.Rapor, error)
	CreateRapor(ctx context.Context, tenantID string, m *model.Rapor) error
	// Tabungan
	GetTabunganList(ctx context.Context, tenantID string) ([]model.Tabungan, error)
	CreateTabunganMutasi(ctx context.Context, tenantID string, m *model.TabunganMutasi) error
	// Kalender
	GetKalenderEvents(ctx context.Context, tenantID string) ([]model.KalenderEvent, error)
	CreateKalenderEvent(ctx context.Context, tenantID string, m *model.KalenderEvent) error
	// Profil
	GetProfil(ctx context.Context, tenantID string) (*model.Profil, error)
	UpdateProfil(ctx context.Context, tenantID string, m *model.ProfilUpdate) error
	// Laporan
	GetReportData(ctx context.Context, tenantID string, req model.ReportRequest) ([]model.ReportData, error)
	// Dashboard Stats
	GetDashboardStats(ctx context.Context, tenantID string) (*model.SekolahDashboardStats, error)
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

// ------ Kelas Implementation ------

func (d *akademikDomain) GetKelasList(ctx context.Context, tenantID string) ([]model.Kelas, error) {
	return d.databasePort.Sekolah().GetKelasByTenant(tenantID)
}

func (d *akademikDomain) CreateKelas(ctx context.Context, tenantID string, kelas *model.Kelas) error {
	kelas.TenantID = tenantID
	return d.databasePort.Sekolah().CreateKelas(kelas)
}

// ------ Asrama Implementation ------

func (d *akademikDomain) GetAsramaList(ctx context.Context, tenantID string) ([]model.Asrama, error) {
	return d.databasePort.Sekolah().GetAsramaByTenant(tenantID)
}

func (d *akademikDomain) CreateAsrama(ctx context.Context, tenantID string, asrama *model.Asrama) error {
	asrama.TenantID = tenantID
	return d.databasePort.Sekolah().CreateAsrama(asrama)
}

func (d *akademikDomain) GetKamarList(ctx context.Context, tenantID, asramaID string) ([]model.Kamar, error) {
	return d.databasePort.Sekolah().GetKamarByAsrama(tenantID, asramaID)
}

func (d *akademikDomain) CreateKamar(ctx context.Context, tenantID string, kamar *model.Kamar) error {
	kamar.TenantID = tenantID
	return d.databasePort.Sekolah().CreateKamar(kamar)
}

func (d *akademikDomain) GetPenempatanList(ctx context.Context, tenantID string) ([]model.Penempatan, error) {
	return d.databasePort.Sekolah().GetPenempatanByTenant(tenantID)
}

func (d *akademikDomain) CreatePenempatan(ctx context.Context, tenantID string, penempatan *model.Penempatan) error {
	penempatan.TenantID = tenantID
	// Optional: Check capacity before placing (skipping for initial version, can be added later)
	return d.databasePort.Sekolah().CreatePenempatan(penempatan)
}

// ------ Kepesantrenan Implementation ------

func (d *akademikDomain) GetPelanggaranAturanList(ctx context.Context, tenantID string) ([]model.PelanggaranAturan, error) {
	return d.databasePort.Sekolah().GetPelanggaranAturan(tenantID)
}

func (d *akademikDomain) CreatePelanggaranAturan(ctx context.Context, tenantID string, m *model.PelanggaranAturan) error {
	m.TenantID = tenantID
	return d.databasePort.Sekolah().CreatePelanggaranAturan(m)
}

func (d *akademikDomain) GetPelanggaranSiswaList(ctx context.Context, tenantID string) ([]model.PelanggaranSiswa, error) {
	return d.databasePort.Sekolah().GetPelanggaranSiswa(tenantID)
}

func (d *akademikDomain) CreatePelanggaranSiswa(ctx context.Context, tenantID string, m *model.PelanggaranSiswa) error {
	m.TenantID = tenantID
	// Optional: Validate AturanID existence and Poin match (skipping for MVP, trusting frontend/db constraints)
	return d.databasePort.Sekolah().CreatePelanggaranSiswa(m)
}

func (d *akademikDomain) GetPerizinanList(ctx context.Context, tenantID string) ([]model.Perizinan, error) {
	return d.databasePort.Sekolah().GetPerizinan(tenantID)
}

func (d *akademikDomain) CreatePerizinan(ctx context.Context, tenantID string, m *model.Perizinan) error {
	m.TenantID = tenantID
	return d.databasePort.Sekolah().CreatePerizinan(m)
}

// ------ Tahfidz Implementation ------

func (d *akademikDomain) GetTahfidzSetoranList(ctx context.Context, tenantID string) ([]model.TahfidzSetoran, error) {
	return d.databasePort.Sekolah().GetTahfidzSetoran(tenantID)
}

func (d *akademikDomain) CreateTahfidzSetoran(ctx context.Context, tenantID string, m *model.TahfidzSetoran) error {
	m.TenantID = tenantID
	return d.databasePort.Sekolah().CreateTahfidzSetoran(m)
}

// ------ Diniyah Implementation ------

func (d *akademikDomain) GetDiniyahKitabList(ctx context.Context, tenantID string) ([]model.DiniyahKitab, error) {
	return d.databasePort.Sekolah().GetDiniyahKitab(tenantID)
}

func (d *akademikDomain) CreateDiniyahKitab(ctx context.Context, tenantID string, m *model.DiniyahKitab) error {
	m.TenantID = tenantID
	return d.databasePort.Sekolah().CreateDiniyahKitab(m)
}

// ------ Rapor Implementation ------

func (d *akademikDomain) GetRaporList(ctx context.Context, tenantID string) ([]model.Rapor, error) {
	return d.databasePort.Sekolah().GetRaporList(tenantID)
}

func (d *akademikDomain) CreateRapor(ctx context.Context, tenantID string, m *model.Rapor) error {
	m.TenantID = tenantID
	return d.databasePort.Sekolah().CreateRapor(m)
}

// ------ Tabungan Implementation ------

func (d *akademikDomain) GetTabunganList(ctx context.Context, tenantID string) ([]model.Tabungan, error) {
	return d.databasePort.Sekolah().GetTabunganList(tenantID)
}

func (d *akademikDomain) CreateTabunganMutasi(ctx context.Context, tenantID string, m *model.TabunganMutasi) error {
	m.TenantID = tenantID
	return d.databasePort.Sekolah().CreateTabunganMutasi(m)
}

// ------ Kalender Implementation ------

func (d *akademikDomain) GetKalenderEvents(ctx context.Context, tenantID string) ([]model.KalenderEvent, error) {
	return d.databasePort.Sekolah().GetKalenderEvents(tenantID)
}

func (d *akademikDomain) CreateKalenderEvent(ctx context.Context, tenantID string, m *model.KalenderEvent) error {
	m.TenantID = tenantID
	return d.databasePort.Sekolah().CreateKalenderEvent(m)
}

// ------ Profil Implementation ------

func (d *akademikDomain) GetProfil(ctx context.Context, tenantID string) (*model.Profil, error) {
	return d.databasePort.Sekolah().GetProfil(tenantID)
}

func (d *akademikDomain) UpdateProfil(ctx context.Context, tenantID string, m *model.ProfilUpdate) error {
	return d.databasePort.Sekolah().UpdateProfil(tenantID, m)
}

// ------ Laporan Implementation ------

func (d *akademikDomain) GetReportData(ctx context.Context, tenantID string, req model.ReportRequest) ([]model.ReportData, error) {
	return d.databasePort.Sekolah().GetReportData(tenantID, req)
}

// ------ Dashboard Stats Implementation ------

func (d *akademikDomain) GetDashboardStats(ctx context.Context, tenantID string) (*model.SekolahDashboardStats, error) {
	return d.databasePort.Sekolah().GetDashboardStats(tenantID)
}
