package domain

import (
	audit_log_domain "eduvera/internal/domain/audit_log"
	"eduvera/internal/domain/auth"
	"eduvera/internal/domain/client"
	"eduvera/internal/domain/content"
	disbursement_domain "eduvera/internal/domain/disbursement"
	erapor_domain "eduvera/internal/domain/erapor"
	notification_domain "eduvera/internal/domain/notification"
	"eduvera/internal/domain/payment"
	dashboard "eduvera/internal/domain/pesantren/dashboard"
	"eduvera/internal/domain/sekolah"
	spp_domain "eduvera/internal/domain/spp"
	"eduvera/internal/domain/tenant"
	outbound_port "eduvera/internal/port/outbound"
)

type Domain interface {
	Client() client.ClientDomain
	Tenant() tenant.TenantDomain
	Auth() auth.AuthDomain
	Payment() payment.PaymentDomain
	Content() content.ContentDomain

	Disbursement() disbursement_domain.Service
	SPP() spp_domain.Service
	PesantrenDashboard() dashboard.Service
	Notification() notification_domain.Service
	Sekolah() sekolah.AkademikDomain
	AuditLog() audit_log_domain.Service
	ERapor() *erapor_domain.Service
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

func (d *domain) Content() content.ContentDomain {
	return content.NewContentDomain(d.databasePort)
}

func (d *domain) Disbursement() disbursement_domain.Service {
	return disbursement_domain.NewService(d.databasePort.Disbursement())
}

func (d *domain) SPP() spp_domain.Service {
	return spp_domain.NewService(d.databasePort.SPP())
}

func (d *domain) PesantrenDashboard() dashboard.Service {
	// Use Mock Service for now
	return dashboard.NewMockService()
}

func (d *domain) Notification() notification_domain.Service {
	return notification_domain.NewService(d.databasePort.Notification())
}

func (d *domain) Sekolah() sekolah.AkademikDomain {
	return sekolah.NewAkademikDomain(d.databasePort)
}

func (d *domain) AuditLog() audit_log_domain.Service {
	return audit_log_domain.NewService(d.databasePort.AuditLog())
}

func (d *domain) ERapor() *erapor_domain.Service {
	return erapor_domain.NewService(d.databasePort.ERapor())
}
