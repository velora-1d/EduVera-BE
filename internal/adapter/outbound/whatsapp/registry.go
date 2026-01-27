package whatsapp_outbound_adapter

import (
	outbound_port "eduvera/internal/port/outbound"
)

type adapter struct {
}

func NewAdapter() outbound_port.MessagePort {
	return &adapter{}
}

func (s *adapter) Client() outbound_port.ClientMessagePort {
	// Not implemented for this adapter
	return nil
}

func (s *adapter) WhatsApp() outbound_port.WhatsAppMessagePort {
	return NewFonnteAdapter()
}
