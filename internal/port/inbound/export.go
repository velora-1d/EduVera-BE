package inbound_port

import "github.com/gofiber/fiber/v2"

type ExportHttpPort interface {
	ExportStudents(c *fiber.Ctx) error
	ExportPayments(c *fiber.Ctx) error
}
