package fiber_inbound_adapter

import (
	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"

	"github.com/gofiber/fiber/v2"
)

type subscriptionAdapter struct {
	domainRegistry domain.Domain
}

func NewSubscriptionAdapter(domainRegistry domain.Domain) inbound_port.SubscriptionHttpPort {
	return &subscriptionAdapter{
		domainRegistry: domainRegistry,
	}
}

// GetSubscription
// GET /api/v1/subscription
func (a *subscriptionAdapter) GetSubscription(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)
	if tenantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tenant_id required"})
	}

	sub, err := a.domainRegistry.Subscription().GetSubscription(c.Context(), tenantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if sub == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Subscription not found"})
	}

	return c.JSON(fiber.Map{"data": sub})
}

// CalculateUpgrade
// POST /api/v1/subscription/calculate-upgrade
func (a *subscriptionAdapter) CalculateUpgrade(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)

	var input model.UpgradeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.TenantID = tenantID // Force tenant ID from token

	calc, err := a.domainRegistry.Subscription().CalculateUpgrade(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": calc})
}

// UpgradePlan
// POST /api/v1/subscription/upgrade
func (a *subscriptionAdapter) UpgradePlan(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)

	var input model.UpgradeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.TenantID = tenantID

	// NOTE: In production, this should probably be called by a payment hook or verify payment first.
	// But per requirements, we allow direct upgrade (maybe assuming free/test or simple flow).
	// Or maybe the frontend calls this AFTER payment success?
	// For now, let's implement validation.

	result, err := a.domainRegistry.Subscription().UpgradePlan(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": result})
}

// DowngradePlan
// POST /api/v1/subscription/downgrade
func (a *subscriptionAdapter) DowngradePlan(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(string)

	var input model.DowngradeInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	input.TenantID = tenantID

	result, err := a.domainRegistry.Subscription().DowngradePlan(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": result})
}
