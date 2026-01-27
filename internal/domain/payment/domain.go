package payment

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/palantir/stacktrace"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

type PaymentDomain interface {
	CreateSnapTransaction(ctx context.Context, input *model.CreatePaymentInput, customerName, customerEmail string) (*model.Payment, *model.SnapTransactionResponse, error)
	HandleWebhook(ctx context.Context, notification *model.MidtransNotification) error
	GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
}

type paymentDomain struct {
	databasePort outbound_port.DatabasePort
	messagePort  outbound_port.MessagePort
	snapClient   snap.Client
}

func NewPaymentDomain(databasePort outbound_port.DatabasePort, messagePort outbound_port.MessagePort) PaymentDomain {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	isProduction := os.Getenv("MIDTRANS_IS_PRODUCTION") == "true"

	var env midtrans.EnvironmentType
	if isProduction {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	var snapClient snap.Client
	snapClient.New(serverKey, env)

	return &paymentDomain{
		databasePort: databasePort,
		messagePort:  messagePort,
		snapClient:   snapClient,
	}
}

func (d *paymentDomain) CreateSnapTransaction(ctx context.Context, input *model.CreatePaymentInput, customerName, customerEmail string) (*model.Payment, *model.SnapTransactionResponse, error) {
	// Calculate price
	amount := model.GetPlanPrice(input.PlanType, input.IsAnnual)
	if amount == 0 {
		return nil, nil, stacktrace.NewError("invalid plan type")
	}

	// Generate order ID
	orderID := model.GenerateOrderID(input.TenantID)

	// Get billing cycle description
	billingCycle := "Monthly"
	if input.IsAnnual {
		billingCycle = "Annual"
	}

	// Create Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    input.PlanType,
				Name:  "EduVera " + input.PlanType + " (" + billingCycle + ")",
				Price: amount,
				Qty:   1,
			},
		},
	}

	// Create Snap transaction
	snapResp, midtransErr := d.snapClient.CreateTransaction(snapReq)
	if midtransErr != nil {
		return nil, nil, stacktrace.Propagate(midtransErr, "failed to create snap transaction")
	}

	// Save payment record
	payment := &model.Payment{
		TenantID:  input.TenantID,
		OrderID:   orderID,
		Amount:    amount,
		Status:    model.PaymentStatusPending,
		SnapToken: snapResp.Token,
	}

	err := d.databasePort.Payment().Create(payment)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "failed to save payment")
	}

	return payment, &model.SnapTransactionResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

func (d *paymentDomain) HandleWebhook(ctx context.Context, notification *model.MidtransNotification) error {
	// Verify signature
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	signatureInput := notification.OrderID + notification.StatusCode + notification.GrossAmount + serverKey
	hash := sha512.Sum512([]byte(signatureInput))
	expectedSignature := hex.EncodeToString(hash[:])

	if notification.SignatureKey != expectedSignature {
		return stacktrace.NewError("invalid signature")
	}

	// Update payment based on transaction status
	var status string
	switch notification.TransactionStatus {
	case "capture", "settlement":
		if notification.FraudStatus == "accept" || notification.FraudStatus == "" {
			status = model.PaymentStatusPaid
			// Mark as paid
			err := d.databasePort.Payment().MarkAsPaid(
				notification.OrderID,
				notification.PaymentType,
				notification.TransactionID,
			)
			if err != nil {
				return stacktrace.Propagate(err, "failed to mark payment as paid")
			}

			// Activate tenant
			payment, err := d.databasePort.Payment().FindByOrderID(notification.OrderID)
			if err == nil && payment != nil {
				_ = d.databasePort.Tenant().UpdateStatus(payment.TenantID, model.TenantStatusActive)

				// Send WhatsApp Notification to Admin
				if d.messagePort != nil {
					// Find admin user for this tenant
					users, err := d.databasePort.User().FindByFilter(model.UserFilter{
						TenantIDs: []string{payment.TenantID},
						Roles:     []string{model.RoleAdmin},
					})
					if err == nil && len(users) > 0 {
						admin := users[0]
						if admin.WhatsApp != "" {
							message := "Halo " + admin.Name + "!\n\n" +
								"Pembayaran Anda untuk Order ID " + notification.OrderID + " telah BERHASIL kami terima.\n\n" +
								"Status: AKTIF\n" +
								"Metode: " + notification.PaymentType + "\n" +
								"Jumlah: Rp " + notification.GrossAmount + "\n\n" +
								"Akun institusi Anda kini telah aktif. Silakan login kembali untuk mulai menggunakan EduVera.\n\n" +
								"Terima kasih atas kepercayaannya!"
							_ = d.messagePort.WhatsApp().Send(admin.WhatsApp, message)
						}
					}
				}
			}

			return nil
		}
	case "pending":
		status = model.PaymentStatusPending
	case "deny", "cancel":
		status = model.PaymentStatusFailed
	case "expire":
		status = model.PaymentStatusExpired
	default:
		status = model.PaymentStatusPending
	}

	return d.databasePort.Payment().UpdateStatus(
		notification.OrderID,
		status,
		notification.PaymentType,
		notification.TransactionID,
	)
}

func (d *paymentDomain) GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	return d.databasePort.Payment().FindByOrderID(orderID)
}
