package fiber_inbound_adapter

import (
	"context"

	"prabogo/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type tenantWhatsAppAdapter struct {
	domainRegistry domain.Domain
}

func NewTenantWhatsAppAdapter(domainRegistry domain.Domain) *tenantWhatsAppAdapter {
	return &tenantWhatsAppAdapter{
		domainRegistry: domainRegistry,
	}
}

// Connect initiates WhatsApp connection for tenant
// POST /api/v1/tenant/whatsapp/connect
func (a *tenantWhatsAppAdapter) Connect(c *fiber.Ctx) error {
	// Get tenant ID from JWT context
	tenantID := c.Locals("tenant_id")
	if tenantID == nil || tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"error":  "Tenant ID not found in token",
		})
	}

	// Check if Premium tier (Premium only feature)
	tier := c.Locals("subscription_tier")
	if tier != "premium" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"error":  "Fitur ini hanya untuk langganan Premium",
		})
	}

	ctx := context.Background()
	session, err := a.domainRegistry.WhatsApp().ConnectTenant(ctx, tenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"qr_code":       session.QRCode,
			"instance_name": session.InstanceName,
			"status":        session.Status,
		},
	})
}

// Status checks WhatsApp connection status
// GET /api/v1/tenant/whatsapp/status
func (a *tenantWhatsAppAdapter) Status(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil || tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"error":  "Tenant ID not found in token",
		})
	}

	ctx := context.Background()
	session, err := a.domainRegistry.WhatsApp().GetStatus(ctx, tenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"connected":     session.Status == "connected",
			"status":        session.Status,
			"instance_name": session.InstanceName,
			"device_info":   session.DeviceInfo,
		},
	})
}

// Disconnect removes WhatsApp connection
// POST /api/v1/tenant/whatsapp/disconnect
func (a *tenantWhatsAppAdapter) Disconnect(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil || tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"error":  "Tenant ID not found in token",
		})
	}

	ctx := context.Background()
	err := a.domainRegistry.WhatsApp().DisconnectTenant(ctx, tenantID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "WhatsApp berhasil diputus",
	})
}

// SendTest sends a test message to verify connection
// POST /api/v1/tenant/whatsapp/test
func (a *tenantWhatsAppAdapter) SendTest(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil || tenantID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"error":  "Tenant ID not found in token",
		})
	}

	var req struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid request body",
		})
	}

	if req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "Phone number is required",
		})
	}

	if req.Message == "" {
		req.Message = "Ini adalah pesan test dari EduVera"
	}

	ctx := context.Background()
	err := a.domainRegistry.WhatsApp().SendMessage(ctx, tenantID.(string), req.Phone, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Pesan test berhasil dikirim",
	})
}
