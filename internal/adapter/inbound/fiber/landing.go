package fiber_inbound_adapter

import (
	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
	inbound_port "prabogo/internal/port/inbound"
)

type landingAdapter struct {
	domain domain.Domain
}

func NewLandingAdapter(domain domain.Domain) inbound_port.LandingHttpPort {
	return &landingAdapter{
		domain: domain,
	}
}

func (h *landingAdapter) Home(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"app":     "EduVera API",
		"version": "1.0.0",
		"status":  "running",
	})
}

// Register (GET) - Removed or keep as metadata provider
func (h *landingAdapter) Register(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"meta": "Register Page Metadata",
	})
}

func (h *landingAdapter) RegisterProcess(a any) error {
	c := a.(*fiber.Ctx)
	// TODO: Port Logic (Input Validation, Service Call)
	// For API, return JSON success
	return c.JSON(fiber.Map{
		"status":    "success",
		"message":   "Registration init success",
		"next_step": "/onboarding/step-2",
		"tenant_id": "temp-123", // Dummy
	})
}

// Step2 (GET) - Metadata
func (h *landingAdapter) Step2(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"options": []string{"School", "Pesantren"},
	})
}

func (h *landingAdapter) Step2Process(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":    "success",
		"next_step": "/onboarding/step-3",
	})
}

// Step3 (GET)
func (h *landingAdapter) Step3(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"meta": "Subdomain Setup",
	})
}

func (h *landingAdapter) Step3Process(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":    "success",
		"next_step": "/onboarding/step-4",
	})
}

// Step4 (GET)
func (h *landingAdapter) Step4(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"plans": []string{"Starter", "Growth", "Enterprise"},
	})
}

func (h *landingAdapter) Step4Process(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":    "success",
		"next_step": "/onboarding/step-5",
	})
}

// Step5 (GET)
func (h *landingAdapter) Step5(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"meta": "Confirmation",
	})
}

func (h *landingAdapter) Step5Process(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":    "success",
		"next_step": "/onboarding/step-6",
	})
}

func (h *landingAdapter) Step6(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":    "pending_payment",
		"va_number": "1234567890",
	})
}

func (h *landingAdapter) Step6Process(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":    "success",
		"next_step": "/onboarding/step-7",
	})
}

func (h *landingAdapter) Step7(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":      "active",
		"permissions": []string{"dashboard_access"},
	})
}

func (h *landingAdapter) Step7Process(a any) error {
	c := a.(*fiber.Ctx)
	return c.JSON(fiber.Map{
		"status":    "success",
		"next_step": "/dashboard",
	})
}
