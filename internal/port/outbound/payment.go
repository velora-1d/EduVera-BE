package outbound_port

import "eduvera/internal/model"

//go:generate mockgen -source=payment.go -destination=./../../../tests/mocks/port/mock_payment.go
type PaymentDatabasePort interface {
	Create(payment *model.Payment) error
	Update(payment *model.Payment) error
	FindByFilter(filter model.PaymentFilter) ([]model.Payment, error)
	FindByID(id string) (*model.Payment, error)
	FindByOrderID(orderID string) (*model.Payment, error)
	UpdateStatus(orderID string, status string, paymentType string, midtransID string) error
	MarkAsPaid(orderID string, paymentType string, midtransID string) error
}
