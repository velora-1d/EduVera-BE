package inbound_port

import "github.com/gofiber/fiber/v2"

type SekolahHttpPort interface {
	GetSiswaList(c *fiber.Ctx) error
	CreateSiswa(c *fiber.Ctx) error
	GetGuruList(c *fiber.Ctx) error
	CreateGuru(c *fiber.Ctx) error
	// Mapel
	GetMapelList(c *fiber.Ctx) error
	GetKelasList(c *fiber.Ctx) error
	CreateKelas(c *fiber.Ctx) error

	// Asrama
	GetAsramaList(c *fiber.Ctx) error
	CreateAsrama(c *fiber.Ctx) error
	GetKamarList(c *fiber.Ctx) error
	CreateKamar(c *fiber.Ctx) error
	GetPenempatanList(c *fiber.Ctx) error
	CreatePenempatan(c *fiber.Ctx) error

	// Kepesantrenan
	GetPelanggaranAturanList(c *fiber.Ctx) error
	CreatePelanggaranAturan(c *fiber.Ctx) error
	GetPelanggaranSiswaList(c *fiber.Ctx) error
	CreatePelanggaranSiswa(c *fiber.Ctx) error
	GetPerizinanList(c *fiber.Ctx) error
	CreatePerizinan(c *fiber.Ctx) error

	// Tahfidz
	GetTahfidzSetoranList(c *fiber.Ctx) error
	CreateTahfidzSetoran(c *fiber.Ctx) error

	// Diniyah
	GetDiniyahKitabList(c *fiber.Ctx) error
	CreateDiniyahKitab(c *fiber.Ctx) error

	// Rapor
	GetRaporList(c *fiber.Ctx) error
	CreateRapor(c *fiber.Ctx) error

	// Tabungan
	GetTabunganList(c *fiber.Ctx) error
	CreateTabunganMutasi(c *fiber.Ctx) error

	// Kalender
	GetKalenderEvents(c *fiber.Ctx) error
	CreateKalenderEvent(c *fiber.Ctx) error

	// Profil
	GetProfil(c *fiber.Ctx) error
	UpdateProfil(c *fiber.Ctx) error

	// Laporan
	GetReportData(c *fiber.Ctx) error
}
