package inbound_port

import "github.com/gofiber/fiber/v2"

// SDMHttpPort defines handlers for SDM module
type SDMHttpPort interface {
	// Employee CRUD
	GetEmployees(c *fiber.Ctx) error
	CreateEmployee(c *fiber.Ctx) error
	UpdateEmployee(c *fiber.Ctx) error
	DeleteEmployee(c *fiber.Ctx) error

	// Payroll
	GetPayrollByPeriod(c *fiber.Ctx) error
	GeneratePayroll(c *fiber.Ctx) error
	MarkPayrollPaid(c *fiber.Ctx) error
	GetPaySlip(c *fiber.Ctx) error
	GetPayrollConfig(c *fiber.Ctx) error
	SavePayrollConfig(c *fiber.Ctx) error

	// Attendance
	GetAttendance(c *fiber.Ctx) error
	RecordAttendance(c *fiber.Ctx) error
	GetAttendanceSummary(c *fiber.Ctx) error
}
