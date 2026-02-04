package outbound_port

import (
	"context"
	"prabogo/internal/model"
)

type PaymentGatewayPort interface {
	CreateInvoice(ctx context.Context, input model.CreatePaymentInput) (*model.Payment, error)
	// VerifyWebhook(payload []byte, signature string) (bool, error) // Can be added later
}
