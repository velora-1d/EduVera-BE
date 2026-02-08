package domain

import (
	analytics_domain "prabogo/internal/domain/analytics"
	audit_log_domain "prabogo/internal/domain/audit_log"
	"prabogo/internal/domain/auth"
	"prabogo/internal/domain/client"
	"prabogo/internal/domain/content"
	disbursement_domain "prabogo/internal/domain/disbursement"
	erapor_domain "prabogo/internal/domain/erapor"
	export_domain "prabogo/internal/domain/export"
	notification_domain "prabogo/internal/domain/notification"
	"prabogo/internal/domain/payment"
	dashboard "prabogo/internal/domain/pesantren/dashboard"
	sdm_domain "prabogo/internal/domain/sdm"
	"prabogo/internal/domain/sekolah"
	spp_domain "prabogo/internal/domain/spp"
	student_domain "prabogo/internal/domain/student"
	"prabogo/internal/domain/subscription"
	"prabogo/internal/domain/tenant"
	whatsapp_domain "prabogo/internal/domain/whatsapp"
	outbound_port "prabogo/internal/port/outbound"
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
	SDM() sdm_domain.SDMDomain
	Subscription() subscription.SubscriptionDomain
	Analytics() analytics_domain.AnalyticsDomain
	Export() export_domain.ExportDomain
	WhatsApp() whatsapp_domain.WhatsAppDomain
	Student() student_domain.StudentDomain
}

type domain struct {
	databasePort     outbound_port.DatabasePort
	messagePort      outbound_port.MessagePort
	cachePort        outbound_port.CachePort
	workflowPort     outbound_port.WorkflowPort
	evolutionPort    outbound_port.WhatsAppClientPort
	notificationPort outbound_port.NotificationServicePort
}

func NewDomain(
	databasePort outbound_port.DatabasePort,
	messagePort outbound_port.MessagePort,
	cachePort outbound_port.CachePort,
	workflowPort outbound_port.WorkflowPort,
	evolutionPort outbound_port.WhatsAppClientPort,
	notificationPort outbound_port.NotificationServicePort,
) Domain {
	return &domain{
		databasePort:     databasePort,
		messagePort:      messagePort,
		cachePort:        cachePort,
		workflowPort:     workflowPort,
		evolutionPort:    evolutionPort,
		notificationPort: notificationPort,
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
	return spp_domain.NewService(d.databasePort.SPP(), d.databasePort.Student(), d.databasePort.User(), d.notificationPort)
}

func (d *domain) PesantrenDashboard() dashboard.Service {
	return dashboard.NewService(d.databasePort.PesantrenDashboard())
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

func (d *domain) SDM() sdm_domain.SDMDomain {
	return sdm_domain.NewSDMDomain(d.databasePort.SDM())
}

func (d *domain) Subscription() subscription.SubscriptionDomain {
	return subscription.NewSubscriptionDomain(d.databasePort.Subscription(), d.databasePort)
}

func (d *domain) Analytics() analytics_domain.AnalyticsDomain {
	return analytics_domain.NewAnalyticsDomain(d.databasePort)
}

func (d *domain) Export() export_domain.ExportDomain {
	return export_domain.NewExportDomain(d.databasePort)
}

func (d *domain) WhatsApp() whatsapp_domain.WhatsAppDomain {
	return whatsapp_domain.NewWhatsAppDomain(d.databasePort, d.evolutionPort)
}

func (d *domain) Student() student_domain.StudentDomain {
	return student_domain.NewStudentDomain(d.databasePort)
}
