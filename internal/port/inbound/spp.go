package inbound_port

import "github.com/gofiber/fiber/v2"

type SPPHttpPort interface {
	List(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	RecordPayment(c *fiber.Ctx) error
	GetStats(c *fiber.Ctx) error
}
