package inbound_port

import "github.com/gofiber/fiber/v2"

type OwnerHttpPort interface {
	Login(c *fiber.Ctx) error
	GetStats(c *fiber.Ctx) error
	GetTenants(c *fiber.Ctx) error
	GetTenantDetail(c *fiber.Ctx) error
	UpdateTenantStatus(c *fiber.Ctx) error

	// Registration logs
	GetRegistrations(c *fiber.Ctx) error

	// SPP Transactions
	GetSPPTransactions(c *fiber.Ctx) error

	// Disbursements
	GetDisbursements(c *fiber.Ctx) error
	ApproveDisbursement(c *fiber.Ctx) error
	RejectDisbursement(c *fiber.Ctx) error

	// Notification logs
	GetNotificationLogs(c *fiber.Ctx) error
}
