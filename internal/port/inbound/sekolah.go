package inbound_port

import "github.com/gofiber/fiber/v2"

type SekolahHttpPort interface {
	GetSiswaList(c *fiber.Ctx) error
	CreateSiswa(c *fiber.Ctx) error
	GetGuruList(c *fiber.Ctx) error
	CreateGuru(c *fiber.Ctx) error
	// Mapel
	GetMapelList(c *fiber.Ctx) error
}
