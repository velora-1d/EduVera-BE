package inbound_port

import "github.com/gofiber/fiber/v2"

// ERaporHttpPort defines handlers for E-Rapor module
type ERaporHttpPort interface {
	// Subject CRUD
	GetSubjects(c *fiber.Ctx) error
	CreateSubject(c *fiber.Ctx) error
	UpdateSubject(c *fiber.Ctx) error
	DeleteSubject(c *fiber.Ctx) error

	// Grade CRUD
	SaveGrade(c *fiber.Ctx) error
	BatchSaveGrades(c *fiber.Ctx) error
	GetStudentGrades(c *fiber.Ctx) error
	GetSubjectGrades(c *fiber.Ctx) error

	// Rapor
	GetStudentRapor(c *fiber.Ctx) error

	// Stats
	GetStats(c *fiber.Ctx) error
}
