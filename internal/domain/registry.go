package domain

import (
	"eduvera/internal/domain/auth"
	"eduvera/internal/domain/client"
	"eduvera/internal/domain/payment"
	"eduvera/internal/domain/tenant"
	outbound_port "eduvera/internal/port/outbound"
)

type Domain interface {
	Client() client.ClientDomain
	Tenant() tenant.TenantDomain
	Auth() auth.AuthDomain
	Payment() payment.PaymentDomain
}

type domain struct {
	databasePort outbound_port.DatabasePort
	messagePort  outbound_port.MessagePort
	cachePort    outbound_port.CachePort
	workflowPort outbound_port.WorkflowPort
}

func NewDomain(
	databasePort outbound_port.DatabasePort,
	messagePort outbound_port.MessagePort,
	cachePort outbound_port.CachePort,
	workflowPort outbound_port.WorkflowPort,
) Domain {
	return &domain{
		databasePort: databasePort,
		messagePort:  messagePort,
		cachePort:    cachePort,
		workflowPort: workflowPort,
	}
}

func (d *domain) Client() client.ClientDomain {
	return client.NewClientDomain(d.databasePort, d.messagePort, d.cachePort, d.workflowPort)
}

func (d *domain) Tenant() tenant.TenantDomain {
	return tenant.NewTenantDomain(d.databasePort)
}

func (d *domain) Auth() auth.AuthDomain {
	return auth.NewAuthDomain(d.databasePort, d.messagePort)
}

func (d *domain) Payment() payment.PaymentDomain {
	return payment.NewPaymentDomain(d.databasePort, d.messagePort)
}
