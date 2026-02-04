package inbound_port

import (
	"github.com/gofiber/fiber/v2"
)

type SubscriptionHttpPort interface {
	GetSubscription(c *fiber.Ctx) error
	GetPricing(c *fiber.Ctx) error
	CalculateUpgrade(c *fiber.Ctx) error
	UpgradePlan(c *fiber.Ctx) error
	DowngradePlan(c *fiber.Ctx) error
}
