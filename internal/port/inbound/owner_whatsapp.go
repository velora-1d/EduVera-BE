package inbound_port

import "github.com/gofiber/fiber/v2"

type OwnerWhatsAppHttpPort interface {
	Connect(c *fiber.Ctx) error
	GetStatus(c *fiber.Ctx) error
	Disconnect(c *fiber.Ctx) error
	TestSend(c *fiber.Ctx) error
}
