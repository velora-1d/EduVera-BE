package fiber_inbound_adapter

import (
	"context"

	"prabogo/internal/domain"
	whatsapp_domain "prabogo/internal/domain/whatsapp"
	"prabogo/internal/model"

	"github.com/gofiber/fiber/v2"
)

type ownerWhatsAppAdapter struct {
	whatsappDomain whatsapp_domain.WhatsAppDomain
}

func NewOwnerWhatsAppAdapter(d domain.Domain) *ownerWhatsAppAdapter {
	return &ownerWhatsAppAdapter{
		whatsappDomain: d.WhatsApp(),
	}
}

// Owner WhatsApp Session - stored separately from tenant sessions
const ownerInstanceName = "eduvera_owner"

// Owner tenant ID placeholder (empty means owner-level session)
const ownerTenantID = ""

// POST /api/v1/owner/whatsapp/connect
func (h *ownerWhatsAppAdapter) Connect(c *fiber.Ctx) error {
	ctx := context.Background()

	// For owner, we use an empty tenant ID which means owner-level session
	session, err := h.whatsappDomain.ConnectTenant(ctx, ownerTenantID)
	if err != nil {
		// Check if already connected
		if err.Error() == "already connected" {
			return c.JSON(fiber.Map{
				"status":  "already_connected",
				"message": "WhatsApp owner sudah terkoneksi",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat sesi WhatsApp: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "pending",
		"qr_code": session.QRCode,
		"message": "Scan QR code dengan WhatsApp",
	})
}

// GET /api/v1/owner/whatsapp/status
func (h *ownerWhatsAppAdapter) GetStatus(c *fiber.Ctx) error {
	ctx := context.Background()

	session, err := h.whatsappDomain.GetStatus(ctx, ownerTenantID)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  "disconnected",
			"message": "WhatsApp belum terkoneksi",
		})
	}

	return c.JSON(fiber.Map{
		"status":       session.Status,
		"phone_number": session.PhoneNumber,
		"message":      getStatusMessage(session.Status),
	})
}

// POST /api/v1/owner/whatsapp/disconnect
func (h *ownerWhatsAppAdapter) Disconnect(c *fiber.Ctx) error {
	ctx := context.Background()

	err := h.whatsappDomain.DisconnectTenant(ctx, ownerTenantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "disconnected",
		"message": "WhatsApp berhasil diputuskan",
	})
}

// POST /api/v1/owner/whatsapp/test
func (h *ownerWhatsAppAdapter) TestSend(c *fiber.Ctx) error {
	ctx := context.Background()

	var input struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if input.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Phone number is required",
		})
	}

	if input.Message == "" {
		input.Message = "Ini adalah pesan test dari EduVera Owner"
	}

	// Send test message using owner's connected WhatsApp
	if err := h.whatsappDomain.SendMessage(ctx, ownerTenantID, input.Phone, input.Message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengirim pesan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "sent",
		"message": "Pesan berhasil dikirim",
	})
}

func getStatusMessage(status string) string {
	switch status {
	case model.WhatsAppStatusConnected:
		return "WhatsApp terkoneksi"
	case model.WhatsAppStatusConnecting:
		return "Menunggu scan QR code"
	default:
		return "WhatsApp tidak terkoneksi"
	}
}
