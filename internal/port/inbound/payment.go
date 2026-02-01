package inbound_port

import "github.com/gofiber/fiber/v2"

//go:generate mockgen -source=payment.go -destination=./../../../tests/mocks/port/mock_payment.go
type PaymentHttpPort interface {
	CreateTransaction(c *fiber.Ctx) error
	Webhook(c *fiber.Ctx) error
	GetStatus(c *fiber.Ctx) error
	// SPP Payment (Premium tier only)
	CreateSPPPayment(c *fiber.Ctx) error
	SPPWebhook(c *fiber.Ctx) error
}
