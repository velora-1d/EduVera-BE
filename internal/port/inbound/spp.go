package inbound_port

import "github.com/gofiber/fiber/v2"

type SPPHttpPort interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	RecordPayment(c *fiber.Ctx) error
	GetStats(c *fiber.Ctx) error
	// Manual payment confirmation methods
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UploadProof(c *fiber.Ctx) error
	ConfirmPayment(c *fiber.Ctx) error
	ListOverdue(c *fiber.Ctx) error
	GenerateManual(c *fiber.Ctx) error
	BroadcastOverdueManual(c *fiber.Ctx) error
}
