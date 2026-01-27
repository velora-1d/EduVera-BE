package outbound_port

//go:generate mockgen -source=whatsapp.go -destination=./../../../tests/mocks/port/mock_whatsapp.go
type WhatsAppMessagePort interface {
	Send(target, message string) error
}
