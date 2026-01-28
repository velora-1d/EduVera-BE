package inbound_port

import "github.com/gofiber/fiber/v2"

type PesantrenDashboardHttpPort interface {
	GetStats(c *fiber.Ctx) error
}
