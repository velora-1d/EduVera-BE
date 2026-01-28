package inbound_port

import "github.com/gofiber/fiber/v2"

type OwnerHttpPort interface {
	Login(c *fiber.Ctx) error
	GetStats(c *fiber.Ctx) error
	GetTenants(c *fiber.Ctx) error
	GetTenantDetail(c *fiber.Ctx) error
	UpdateTenantStatus(c *fiber.Ctx) error
}
