package fiber_inbound_adapter

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
	outbound_port "prabogo/internal/port/outbound"
)

type paymentAdapter struct {
	domainRegistry domain.Domain
	whatsapp       outbound_port.WhatsAppMessagePort
}

func NewPaymentAdapter(
	domainRegistry domain.Domain,
	whatsapp outbound_port.WhatsAppMessagePort,
) inbound_port.PaymentHttpPort {
	return &paymentAdapter{
		domainRegistry: domainRegistry,
		whatsapp:       whatsapp,
	}
}

type CreatePaymentRequest struct {
	TenantID     string `json:"tenant_id"`
	PlanType     string `json:"plan_type"`
	IsAnnual     bool   `json:"is_annual"`
	CustomerName string `json:"customer_name"`
	Email        string `json:"email"`
}

// CreateTransaction creates Midtrans Snap transaction
// POST /api/v1/payment/create
func (a *paymentAdapter) CreateTransaction(c *fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.TenantID == "" || req.PlanType == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "tenant_id, plan_type, and email are required",
		})
	}

	ctx := context.Background()
	input := &model.CreatePaymentInput{
		TenantID: req.TenantID,
		PlanType: req.PlanType,
		IsAnnual: req.IsAnnual,
	}

	payment, snapResp, err := a.domainRegistry.Payment().CreateSnapTransaction(
		ctx, input, req.CustomerName, req.Email,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":       "success",
		"order_id":     payment.OrderID,
		"snap_token":   snapResp.Token,
		"redirect_url": snapResp.RedirectURL,
		"amount":       payment.Amount,
	})
}

// Webhook handles Midtrans payment notification
// POST /api/v1/payment/webhook
func (a *paymentAdapter) Webhook(c *fiber.Ctx) error {
	var notification model.MidtransNotification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification payload",
		})
	}

	ctx := context.Background()
	if err := a.domainRegistry.Payment().HandleWebhook(ctx, &notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Handle subscription renewal if payment success
	if notification.TransactionStatus == "settlement" || notification.TransactionStatus == "capture" {
		// Get payment details to find TenantID
		payment, err := a.domainRegistry.Payment().GetPaymentByOrderID(ctx, notification.OrderID)
		if err == nil && payment != nil {
			// Renew subscription
			if err := a.domainRegistry.Subscription().RenewSubscription(ctx, payment.TenantID, notification.OrderID); err != nil {
				// Log error but don't fail the webhook response
			}

			// Get tenant info for notification
			tenant, tenantErr := a.domainRegistry.Tenant().FindByID(ctx, payment.TenantID)
			if tenantErr == nil && tenant != nil && a.whatsapp != nil {
				// Get admin user phone from tenant (find owner)
				users, _ := a.domainRegistry.Auth().GetCurrentUser(ctx, "")
				// For now, use tenant-level notification if we have contact
				// TODO: Get owner phone from users table

				// Send WhatsApp notification (use payment email as fallback lookup)
				if users != nil && users.WhatsApp != "" {
					go func() {
						loginURL := "https://" + tenant.Subdomain + ".eduvera.ve-lora.my.id/login"
						amountStr := fmt.Sprintf("Rp %d", payment.Amount)

						message := fmt.Sprintf(
							"✅ *Pembayaran Berhasil!*\n\n"+
								"Terima kasih! Pembayaran untuk *%s* telah kami terima.\n\n"+
								"💰 *Total:* %s\n"+
								"📋 *Order ID:* %s\n\n"+
								"Langganan Anda sudah aktif!\n\n"+
								"📌 *Dashboard:*\n%s\n\n"+
								"_Tim EduVera_",
							tenant.Name,
							amountStr,
							payment.OrderID,
							loginURL,
						)
						_ = a.whatsapp.Send(users.WhatsApp, message)
					}()
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

// GetStatus returns payment status by order ID
// GET /api/v1/payment/status/:order_id
func (a *paymentAdapter) GetStatus(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order_id is required",
		})
	}

	ctx := context.Background()
	payment, err := a.domainRegistry.Payment().GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Payment not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":         "success",
		"order_id":       payment.OrderID,
		"amount":         payment.Amount,
		"payment_status": payment.Status,
		"payment_type":   payment.PaymentType,
		"paid_at":        payment.PaidAt,
	})
}
