package inbound_port

import "github.com/gofiber/fiber/v2"

//go:generate mockgen -source=whatsapp.go -destination=./../../../tests/mocks/port/mock_whatsapp.go
type TenantWhatsAppHttpPort interface {
	Connect(c *fiber.Ctx) error    // POST /tenant/whatsapp/connect
	Status(c *fiber.Ctx) error     // GET /tenant/whatsapp/status
	Disconnect(c *fiber.Ctx) error // POST /tenant/whatsapp/disconnect
	SendTest(c *fiber.Ctx) error   // POST /tenant/whatsapp/test
}
