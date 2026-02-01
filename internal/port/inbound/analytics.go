package inbound_port

import "github.com/gofiber/fiber/v2"

type AnalyticsHttpPort interface {
	GetAnalytics(c *fiber.Ctx) error
}
