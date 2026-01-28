package inbound_port

import "github.com/gofiber/fiber/v2"

type ContentHttpPort interface {
	Get(c *fiber.Ctx) error
	Upsert(c *fiber.Ctx) error
}
