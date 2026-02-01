package fiber_inbound_adapter

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// SECURITY: Verify Midtrans signature to prevent fake webhooks
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if !notification.VerifySignature(serverKey) {
		log.Printf("[SECURITY] Invalid webhook signature for order: %s", notification.OrderID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid signature",
		})
	}

	ctx := context.Background()
	if err := a.domainRegistry.Payment().HandleWebhook(ctx, &notification); err != nil {
		// Log internal error but return generic message
		log.Printf("[ERROR] Webhook processing failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to process notification",
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

// CreateSPPPaymentRequest for SPP payment creation
type CreateSPPPaymentRequest struct {
	SPPID       string `json:"spp_id"`
	TenantID    string `json:"tenant_id"`
	Amount      int64  `json:"amount"`
	StudentName string `json:"student_name"`
	ParentEmail string `json:"parent_email"`
}

// CreateSPPPayment creates Midtrans Snap for SPP (Premium tier only)
// POST /api/v1/payment/spp/create
func (a *paymentAdapter) CreateSPPPayment(c *fiber.Ctx) error {
	var req CreateSPPPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.SPPID == "" || req.TenantID == "" || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "spp_id, tenant_id, dan amount wajib diisi",
		})
	}

	// Check tenant tier - only Premium can use PG
	ctx := context.Background()
	tenant, err := a.domainRegistry.Tenant().FindByID(ctx, req.TenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tenant tidak ditemukan",
		})
	}

	if !model.HasFeature(tenant.SubscriptionTier, model.FeaturePaymentGateway) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":       "Fitur Payment Gateway hanya tersedia untuk paket Premium",
			"upgrade_url": "/pricing",
		})
	}

	// Create Midtrans Snap transaction for SPP
	payment, snapResp, err := a.domainRegistry.Payment().CreateSPPSnapTransaction(
		ctx, req.SPPID, req.TenantID, req.Amount, req.StudentName, req.ParentEmail,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat transaksi pembayaran. " + err.Error(),
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

// SPPWebhook handles Midtrans callback for SPP payments
// POST /api/v1/payment/spp/webhook
func (a *paymentAdapter) SPPWebhook(c *fiber.Ctx) error {
	var notification model.MidtransNotification
	if err := c.BodyParser(&notification); err != nil {
		log.Printf("[ERROR] SPP Webhook: Failed to parse notification: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification payload",
		})
	}

	// SECURITY: Verify Midtrans signature
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if !notification.VerifySignature(serverKey) {
		log.Printf("[SECURITY] Invalid SPP webhook signature for order: %s", notification.OrderID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid signature",
		})
	}

	log.Printf("[INFO] SPP Webhook received: OrderID=%s, Status=%s", notification.OrderID, notification.TransactionStatus)

	ctx := context.Background()
	if err := a.domainRegistry.Payment().HandleSPPWebhook(ctx, &notification); err != nil {
		log.Printf("[ERROR] SPP Webhook: Error handling notification: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process notification",
		})
	}

	// If payment successful, update SPP status
	if notification.TransactionStatus == "settlement" || notification.TransactionStatus == "capture" {
		// Extract SPP ID from order ID (format: SPP-{sppID}-{timestamp})
		parts := splitOrderID(notification.OrderID)
		if len(parts) >= 2 {
			sppID := parts[1]
			// Update SPP status to paid
			_ = a.domainRegistry.SPP().ConfirmPayment(ctx, sppID, "system", "midtrans")
		}
	}

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

// splitOrderID splits SPP order ID into parts
func splitOrderID(orderID string) []string {
	var parts []string
	current := ""
	for _, c := range orderID {
		if c == '-' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
